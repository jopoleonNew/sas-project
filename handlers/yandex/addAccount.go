package yandex

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	api "github.com/nk2ge5k/goyad"
	"github.com/nk2ge5k/goyad/agencyclients"
	"github.com/nk2ge5k/goyad/campaigns"
	"github.com/nk2ge5k/goyad/clients"
	"github.com/nk2ge5k/goyad/gc"
	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

// GetYandexAuthLink writes to ResponseWriter the Yandex.Direct API Auth Link
// which front-end uses to redirect client to give access to his Yandex.Direct account
func GetYandexAuthLink(w http.ResponseWriter, r *http.Request) {
	accountlogin := r.FormValue("accountlogin")
	log.Println(".............GetYandexAuthLink used: ", accountlogin, Config.YandexDirectAppID)
	yandexUrl := "https://oauth.yandex.ru/authorize?response_type=code&client_id=" + Config.YandexDirectAppID + "&state=" + accountlogin + "&login_hint=" + accountlogin + "&force_confirm=yes"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Write([]byte(yandexUrl))
	return
}

var logr = logrus.New()

func GetYandexToken(w http.ResponseWriter, r *http.Request) {
	logr.Info("GetYandexToken used with request: ", r)
	query := r.URL.Query()
	logr.Info("GetYandexToken income URL query: ", r.URL.Query())

	s := query["code"]
	if s == nil || len(s) == 0 {
		logr.Error("Request from YandexOuath received without code")
		http.Error(w, fmt.Sprintf("Request from Yandex received without code %s", query), http.StatusBadRequest)
		return
	}
	// "state" is the Yandex account login sent with GetAuthCodeYandexHandler()
	al := query["state"]
	if al == nil || len(al) == 0 {
		logr.Error("Request from Yandex received without state")
		http.Error(w, fmt.Sprintf("Request from Yandex received without state %s", query), http.StatusBadRequest)
		return
	}
	accountlogin := al[0]
	code := s[0]
	oauthresp, err := yad.MakeYandexOauthRequest(code, yad.API_YANDEX_OAUTH_URL)
	if err != nil {
		logr.Println("GetYandexAccessToken MakeYandexOauthRequest error: ", err)
		http.Error(w, fmt.Sprintf("cant get auth token from yandex with code: %s, error: %v", code, err), http.StatusBadRequest)
		return
	}
	creator := r.Context().Value("username").(string)
	if creator == "" {
		logr.Errorf("GetYandexToken r.Context().Value(username) is empty: ", creator)
		http.Error(w, fmt.Sprintf("Can't identify username inside request context: %s", creator), http.StatusBadRequest)
		return
	}
	client := api.NewClient()
	client.Token = api.Token{Value: oauthresp.AccessToken}
	client.Login = accountlogin
	_, err = CollectAccountandAddtoBD(client, creator)
	if err != nil {
		logr.Errorf("CollectAccountandAddtoBD error: %v", err)
		http.Error(w, fmt.Sprintf("CollectAccountandAddtoBD %s, %s error: %v", client, creator, err), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)

}

type CreateInfo struct {
	Status          string
	Role            string
	CampaignsAmount int
}

func CollectAccountandAddtoBD(client api.Client, creator string) (info CreateInfo, err error) {
	resultCamps, err := collectCampaings(client)
	if err != nil {
		if strings.Contains(err.Error(), "53") {
			//if that error occurs, this means, that the new account user trying
			// to add is agency, so we must collect agency's clienst and then
			// add new account in DB for each of agency client
			result, err := addYandexAgencyAccounts(client, creator)
			if err != nil {
				logr.Errorf("CollectAccountandAddtoBD addYandexAgencyAccounts(%s, %s, %s) error: %v", client.Login, client.Token.GetToken(), creator, err)
				return info, fmt.Errorf("cant recieve campaings from Yandex.Direct with parametrs %s %s error: ", client.Login, creator, err)

			}
			logrus.Info("addYandexAgencyAccounts SUCCESS")
			_ = result
			return info, nil
		} else {
			logr.Errorln("CollectAccountandAddtoBD unknow error: ", err)
			return info, fmt.Errorf("cant recieve campaings from Yandex.Direct inside CollectAccountandAddtoBD function, error: %v", err)
		}
		return
	}
	resultClientInfo, err := collectClientInfo(client)
	if err != nil {
		logr.Println("GetYandexAccessToken MakeYandexOauthRequest error: ", err)
		return info, fmt.Errorf("cant collectClientInfo from yandex with login: %s, error: %v", client.Login, err)
	}
	logrus.Infof("Information about Client: %+v", resultClientInfo)
	//for _, cli := range resultClientInfo.Clients {
	//	logrus.Infof("Information about Client: %+v", cli)
	//}
	a := model.NewAccount2(
		creator,
		"Яндекс Директ",
		resultClientInfo.Clients[0].Login,
		resultClientInfo.Clients[0].Representatives[0].Email,
	)
	a.Role = "client"
	a.Status = "active"
	a.Owners = append([]string{}, creator)
	a.AccountType = resultClientInfo.Clients[0].Type
	a.AuthToken = client.Token.GetToken()
	a.AppID = Config.YandexDirectAppID
	a.AppSecret = Config.YandexDirectAppSecret
	a.CampaignsInfo = model.AdaptYandexCampaings(resultCamps)
	a.CreatedAt = time.Now()
	err = a.AdvanceUpdate()
	if err != nil {
		logr.Errorf("cant a.AdvanceUpdate() to DB %v \n error: %v", a.Accountlogin, err)
		return info, fmt.Errorf("cant add account to DB %v \n error: %v", a.Accountlogin, err)
	}
	return info, nil
}

func addYandexAgencyAccounts(client api.Client, creator string) (info CreateInfo, err error) {

	var YandexConnectionsLimit = 5
	chAC := make(chan gc.ClientGetItem, 4) // This number 3 can be anything as long as it's larger than YandexConnectionsLimit
	var wg sync.WaitGroup
	resultA, err := collectAgencyClients(client)
	if err != nil {
		logr.Errorln("collectCampaingsfromAgency  error: ", err)
		return
	}
	logrus.Debug("Debug 2 ")
	//var camps []campaigns.GetResponse
	logrus.Debug("Debug 3 ")
	for i := 0; i < YandexConnectionsLimit; i++ {
		wg.Add(1)
		go func() {
			for {
				ac, ok := <-chAC
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				for _, c := range ac.Representatives {
					//collecting ads campaigns from Yandex
					ci := api.NewClient()
					ci.Login = c.Login
					ci.Token = api.Token{Value: client.Token.GetToken()}
					result, err := collectCampaings(ci)
					if err != nil {
						logr.Errorf("cant collect campaings with parameters: collectCampaings(%s, %s) error: %v", ci.Login, ci.Token.GetToken(), err)
					}
					a := model.NewAccount2(
						creator,
						"Яндекс Директ",
						c.Login,
						c.Email,
					)
					a.Role = "client"
					a.Status = "active"
					a.Owners = append([]string{}, creator)
					a.AccountType = ac.Type
					a.AuthToken = client.Token.GetToken()
					a.AppID = Config.YandexDirectAppID
					a.AppSecret = Config.YandexDirectAppSecret
					a.CampaignsInfo = model.AdaptYandexCampaings(result)
					a.CreatedAt = time.Now()
					err = a.AdvanceUpdate()
					if err != nil {
						logr.Errorf("addYandexAgencyAccounts a.AdvanceUpdate(%s) error: %v", c.Login, err)
						return
					}
				}
			}
		}()
	}
	var agencyClients []string
	for _, c := range resultA.Clients {
		logrus.Infof("for resultA.Clients c: %v", c)
		for _, info := range c.Representatives {
			agencyClients = append(agencyClients, info.Login)
		}
		chAC <- c // add client to the queue
	}
	a := model.NewAccount2(
		creator,
		"Яндекс Директ",
		client.Login,
		client.Login+"@yandex.ru",
	)
	a.Creator = creator
	a.Source = "Яндекс Директ"
	a.Accountlogin = client.Login
	a.Email = client.Login + "@yandex.ru"
	a.Role = "agency"
	a.Status = "active"
	a.AgencyClients = agencyClients
	a.Owners = []string{creator}
	a.AuthToken = client.Token.GetToken()
	a.AppID = Config.YandexDirectAppID
	a.AppSecret = Config.YandexDirectAppSecret
	a.CreatedAt = time.Now()
	err = a.AdvanceUpdate()
	if err != nil {
		logr.Errorf("cant add account to DB %v \n error: %v", client.Login, err)
		return
	}
	close(chAC) // This tells the goroutines there's nothing else to do
	wg.Wait()   // Wait for the threads to finish
	return info, nil
}

func collectAgencyClients(client api.Client) (res agencyclients.GetResponse, err error) {
	//logrus.Debug("Debug 1 collectAgencyClients start")
	//client := api.NewClient()
	//client.Token = api.Token{Value: token}
	//client.Login = login
	clientInfo := agencyclients.GetRequest{
		FieldNames: []agencyclients.AgencyClientFieldEnum{
			"AccountQuality", "ClientId", "ClientInfo", "Login", "Phone", "Representatives", "Restrictions", "Type",
		},
	}
	service2 := agencyclients.New(&client)
	result, err := service2.Get(clientInfo)
	if err != nil {
		return res, fmt.Errorf("collectAgencyClients service.Get error %v", err)
	}
	//logrus.Debug("Debug 1 collectAgencyClients end")
	return result, nil
}

func collectClientInfo(client api.Client) (res clients.GetResponse, err error) {
	//client := api.NewClient()
	////clients.New()"f20.ru", "API_TOKEN"
	//client.Token = api.Token{Value: token}
	//client.Login = login
	clientInfo := clients.GetRequest{
		FieldNames: []clients.ClientFieldEnum{
			"ClientId", "ClientInfo", "CountryId", "CreatedAt", "Login", "Representatives", "Type",
		},
	}
	service2 := clients.New(&client)
	result, err := service2.Get(clientInfo)
	if err != nil {
		return res, fmt.Errorf("collectClientInfo service.Get error %v", err)
	}
	return result, nil
}

func collectCampaings(client api.Client) (res campaigns.GetResponse, err error) {

	//DRAFT Кампания создана и еще не отправлена на модерацию.
	//MODERATION	Кампания находится на модерации.
	//ACCEPTED	Хотя бы одно объявление в кампании принято модерацией.
	//REJECTED	Все объявления в кампании отклонены модерацией.
	//UNKNOWN	Используется для обеспечения обратной совместимости и отображения статусов, не поддерживаемых в данной версии API.
	campaingsInfo := campaigns.GetRequest{
		FieldNames: []campaigns.CampaignFieldEnum{
			"ClientInfo",
			"Id",
			"Name",
			"RepresentedBy",
			"Statistics",
			"Status",
			"Type",
		},
	}

	service := campaigns.New(&client)

	result, err := service.Get(campaingsInfo)
	if err != nil {
		return res, fmt.Errorf("collectCampaings service.Get error %v", err)
	}

	return result, nil
}
