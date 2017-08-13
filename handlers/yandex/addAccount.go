package yandex

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nk2ge5k/goyad"
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

// AddYandexAccount
func AddYandexAccount(w http.ResponseWriter, r *http.Request) {
	logrus.Info("AddYandexAccount used with request: ", r)
	query := r.URL.Query()
	logrus.Info("AddYandexAccount income URL query: ", r.URL.Query())

	s := query["code"]
	if s == nil || len(s) == 0 {
		logrus.Error("Request from YandexOuath received without code")
		http.Error(w, fmt.Sprintf("Request from Yandex received without code %s", query), http.StatusBadRequest)
		return
	}
	// "state" is the Yandex account login sent with GetAuthCodeYandexHandler()
	al := query["state"]
	if al == nil || len(al) == 0 {
		logrus.Error("Request from Yandex received without state")
		http.Error(w, fmt.Sprintf("Request from Yandex received without state %s", query), http.StatusBadRequest)
		return
	}
	accountlogin := al[0]
	code := s[0]

	//for testing through request context
	var authURL string
	if r.Context().Value("authurl") == nil {
		authURL = yad.API_YANDEX_OAUTH_URL
	} else {
		authURL = r.Context().Value("authurl").(string)
	}
	oauthresp, err := yad.GetYandexToken(code, authURL)
	if err != nil {
		logrus.Println("GetYandexAccessToken GetYandexToken error: ", err)
		http.Error(w, fmt.Sprintf("cant get auth token from yandex with code: %s, error: %v", code, err), http.StatusBadRequest)
		return
	}
	creator := r.Context().Value("username").(string)
	if creator == "" {
		logrus.Errorf("AddYandexAccount r.Context().Value(username) is empty: ", creator)
		http.Error(w, fmt.Sprintf("Can't identify username inside request context: %s", creator), http.StatusBadRequest)
		return
	}
	client := goyad.NewClient()
	client.Token = goyad.Token{Value: oauthresp.AccessToken}
	client.Login = accountlogin
	//for testing
	//apiURL := r.Context().Value("apiurl").(string)
	//if apiURL != "" {
	//	client.ApiUrl = apiURL
	//}
	_, err = CollectAccountandAddtoBD(client, creator)
	if err != nil {
		logrus.Errorf("CollectAccountandAddtoBD error: %v", err)
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

func CollectAccountandAddtoBD(client goyad.Client, creator string) (info CreateInfo, err error) {
	resultCamps, err := collectCampaings(client)
	if err != nil {
		//if that error occurs, this means, that the new account that user trying
		// to add is agency, so we must collect agency's clients and then
		// add new account in DB for each of agency client
		if strings.Contains(err.Error(), "53") {
			_, err := addYandexAgencyAccounts(client, creator)
			if err != nil {
				logrus.Errorf("CollectAccountandAddtoBD addYandexAgencyAccounts(%s, %s, %s) error: %v", client.Login, client.Token.GetToken(), creator, err)
				return info, fmt.Errorf("cant recieve campaings from Yandex.Direct with parametrs %s %s error: ", client.Login, creator, err)

			}
			logrus.Info("addYandexAgencyAccounts SUCCESS")
			return info, nil
		} else {
			logrus.Errorln("CollectAccountandAddtoBD unknow error: ", err)
			return info, fmt.Errorf("cant recieve campaings from Yandex.Direct inside CollectAccountandAddtoBD function, error: %v", err)
		}
		return
	}
	resultClientInfo, err := collectClientInfo(client)
	if err != nil {
		logrus.Println("GetYandexAccessToken GetYandexToken error: ", err)
		return info, fmt.Errorf("cant collectClientInfo from yandex with login: %s, error: %v", client.Login, err)
	}
	logrus.Infof("Information about Client: %+v", resultClientInfo)
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
	err = a.Update()
	if err != nil {
		logrus.Errorf("cant a.Update() to DB %v \n error: %v", a.Accountlogin, err)
		return info, fmt.Errorf("cant add account to DB %v \n error: %v", a.Accountlogin, err)
	}
	return info, nil
}

func addYandexAgencyAccounts(client goyad.Client, creator string) (info CreateInfo, err error) {

	var YandexConnectionsLimit = 5
	chAC := make(chan gc.ClientGetItem, 4) // channel's buffer is the number of simultaneous gorouitenes
	var wg sync.WaitGroup
	resultA, err := collectAgencyClients(client)
	if err != nil {
		logrus.Errorln("collectCampaingsfromAgency  error: ", err)
		return info, err
	}
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
					//collecting ads campaigns from Yandex for every agency client
					ci := goyad.NewClient()
					ci.Login = c.Login
					ci.Token = goyad.Token{Value: client.Token.GetToken()}
					result, err := collectCampaings(ci)
					if err != nil {
						logrus.Errorf("cant collect campaings with parameters: collectCampaings(%s, %s) error: %v", ci.Login, ci.Token.GetToken(), err)
					}

					//creating new account for each of agency clients
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
					err = a.Update()
					if err != nil {
						logrus.Errorf("addYandexAgencyAccounts a.Update(%s) error: %v", c.Login, err)
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
	err = a.Update()
	if err != nil {
		logrus.Errorf("cant add account to DB %v \n error: %v", client.Login, err)
		return info, err
	}
	close(chAC) // This tells the goroutines there's nothing else to do
	wg.Wait()   // Wait for the threads to finish
	return info, nil
}

func collectAgencyClients(client goyad.Client) (res agencyclients.GetResponse, err error) {
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
	return result, nil
}

func collectClientInfo(client goyad.Client) (res clients.GetResponse, err error) {
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

func collectCampaings(client goyad.Client) (res campaigns.GetResponse, err error) {

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
