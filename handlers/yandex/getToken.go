package yandex

import (
	//"log"
	"net/http"

	"fmt"

	"strings"

	"time"

	"sync"

	api "github.com/nk2ge5k/goyad"
	"github.com/nk2ge5k/goyad/agencyclients"
	"github.com/nk2ge5k/goyad/campaigns"
	"github.com/nk2ge5k/goyad/clients"
	"github.com/nk2ge5k/goyad/gc"
	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
	yad "gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

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
	oauthresp, err := yad.MakeYandexOauthRequest(code)
	if err != nil {
		logr.Println("GetYandexAccessToken MakeYandexOauthRequest error: ", err)
		http.Error(w, fmt.Sprintf("cant get auth token from yandex with code: %s, error: %v", code, err), http.StatusBadRequest)
		return
	}
	resultCamps, err := collectCampaings(accountlogin, oauthresp.AccessToken)
	if err != nil {
		if strings.Contains(err.Error(), "53") {
			result, err := AddNewYandexAccount(accountlogin, oauthresp.AccessToken, r.Context().Value("username").(string))
			if err != nil {
				logr.Errorf("GetYandexAccessToken AddNewYandexAccount(%s, %s, %s) error: %v", accountlogin, oauthresp.AccessToken, r.Context().Value("username").(string), err)
				http.Error(w, fmt.Sprintf("cant recieve campaings from Yandex.Direct with parametrs %s %s error: ", accountlogin, r.Context().Value("username").(string), err), http.StatusBadRequest)
				return
			}
			logrus.Info("AddNewYandexAccount SUCCESS")
			_ = result

		} else {
			logr.Errorln("GetYandexAccessToken collectCampaings error: ", err)
			http.Error(w, fmt.Sprintf("cant recieve campaings from Yandex.Direct %s", err), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return

	}
	resultClientInfo, err := collectClientInfo(accountlogin, oauthresp.AccessToken)
	if err != nil {
		logr.Println("GetYandexAccessToken MakeYandexOauthRequest error: ", err)
		http.Error(w, fmt.Sprintf("cant collectClientInfo from yandex with login: %s, error: %v", accountlogin, err), http.StatusBadRequest)
		return
	}
	for _, cli := range resultClientInfo.Clients {
		logrus.Infof("Information about Client: %+v", cli)
	}

	//logrus.Printf("Result from Yandex about %s: \n %+v", accountlogin, result.Campaigns[5])
	//logrus.Infof("Request Context: %+v", r.Context())
	newAccount := model.Account2{
		Creator:       r.Context().Value("username").(string),
		Source:        "Яндекс Директ",
		Accountlogin:  resultClientInfo.Clients[0].Login,
		Email:         resultClientInfo.Clients[0].Representatives[0].Email,
		Role:          "client",
		Status:        "active",
		AccountType:   resultClientInfo.Clients[0].Type,
		AuthToken:     oauthresp.AccessToken,
		AppID:         Config.YandexDirectAppID,
		AppSecret:     Config.YandexDirectAppSecret,
		CampaignsInfo: model.AdaptYandexCampaings(resultCamps),
		CreatedAt:     time.Now(),
	}
	//acccamps := model.AdaptYandexCampaings(resultCamps)
	//account := model.NewAccount2("", "Яндекс Директ", resultClientInfo.Clients[0].Login, resultClientInfo.Clients[0].Representatives[0].Email)
	//account.CampaignsInfo = acccamps
	//account.Role = "client"
	//account.AuthToken = oauthresp.AccessToken
	//account.AppID = Config.YandexDirectAppID
	//account.AppSecret = Config.YandexDirectAppSecret
	//account.Owners = []string{r.Context().Value("username").(string)}
	//account.CreatedAt = time.Now()
	err = newAccount.AdvanceUpdate()
	if err != nil {
		logr.Errorf("cant add account to DB %v \n error: %v", newAccount.Accountlogin, err)
		return
	}
	//logrus.Printf("Result from Yandex after adaption: %+v", acccamps)

	http.Redirect(w, r, "/accounts", http.StatusSeeOther)

}

type CreateInfo struct {
	Status          string
	Role            string
	CampaignsAmount int
}

func AddNewYandexAccount(login, token, creator string) (info CreateInfo, err error) {

	var YandexConnectionsLimit = 5
	chAC := make(chan gc.ClientGetItem, 4) // This number 3 can be anything as long as it's larger than YandexConnectionsLimit
	var wg sync.WaitGroup
	resultA, err := collectAgencyClients(login, token)
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
					result, err := collectCampaings(c.Login, token)
					if err != nil {
						logr.Errorf("cant collect campaings with parameters: collectCampaings(%s, %s) error: %v", login, token, err)
					}
					account := model.NewAccount2("", "Яндекс Директ", c.Login, c.Email)
					account.CampaignsInfo = model.AdaptYandexCampaings(result)
					account.Role = "client"
					account.AuthToken = token
					account.Creator = creator
					account.AppID = Config.YandexDirectAppID
					account.AppSecret = Config.YandexDirectAppSecret
					account.CreatedAt = time.Now()
					err = account.AdvanceUpdate()
					if err != nil {
						logr.Errorf("cant add account to DB %v \n error: %v", account, err)
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
	account := model.NewAccount2(creator, "Яндекс Директ", login, login+"@yandex.ru")
	account.Role = "agency"
	account.AgencyClients = agencyClients
	account.AuthToken = token
	account.AppID = Config.YandexDirectAppID
	account.AppSecret = Config.YandexDirectAppSecret
	account.Owners = []string{creator}
	account.CreatedAt = time.Now()
	err = account.AdvanceUpdate()
	if err != nil {
		logr.Errorf("cant add account to DB %v \n error: %v", account.Accountlogin, err)
		return
	}
	close(chAC) // This tells the goroutines there's nothing else to do
	wg.Wait()   // Wait for the threads to finish
	return info, nil
}

func collectAgencyClients(login, token string) (res agencyclients.GetResponse, err error) {
	logrus.Debug("Debug 1 collectAgencyClients start")
	client := api.NewClient()
	client.Token = api.Token{Value: token}
	client.Login = login
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
	logrus.Debug("Debug 1 collectAgencyClients end")
	return result, nil
}

func collectCampaingsfromAgency(login, token string) (res []campaigns.GetResponse, err error) {
	fmt.Println("collectCampaingsfromAgency START", time.Now())
	resultA, err := collectAgencyClients(login, token)
	if err != nil {
		logr.Errorln("collectCampaingsfromAgency  error: ", err)
		return
	}
	var camps []campaigns.GetResponse
	for _, ag := range resultA.Clients {
		for _, login := range ag.Representatives {
			result, err := collectCampaings(login.Login, token)
			if err != nil {
				logr.Errorln("cant collect campaings with parameters: collectCampaings(%s, %s) error: %v", login.Login, token, err)
				return nil, fmt.Errorf("cant collect campaings with parameters: collectCampaings(%s, %s), error: %v", login.Login, token, err)
			}
			camps = append(camps, result)
		}
	}
	fmt.Println("collectCampaingsfromAgency END", time.Now())
	return camps, nil
}
func collectCampaingsfromAgencyConcurently(login, token string) (res []campaigns.GetResponse, err error) {
	fmt.Println("collectCampaingsfromAgencyConcurently START", time.Now())

	var YandexConnectionsLimit = 5
	chAC := make(chan gc.ClientGetItem, 4) // This number 3 can be anything as long as it's larger than YandexConnectionsLimit
	var wg sync.WaitGroup
	resultA, err := collectAgencyClients(login, token)
	if err != nil {
		logr.Errorln("collectCampaingsfromAgency  error: ", err)
		return
	}
	logrus.Debug("Debug 2 ")
	var camps []campaigns.GetResponse
	logrus.Debug("Debug 3 ")
	for i := 0; i < YandexConnectionsLimit; i++ {
		logrus.Info("Iterating through YandexConnectionsLimit")
		wg.Add(1)
		go func() {
			for {
				login, ok := <-chAC
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				logrus.Info("Inside gorouitne with login:", login)
				result, err := collectCampaings(login.Login, token)
				if err != nil {
					logr.Errorln("cant collect campaings with parameters: collectCampaings(%s, %s) error: %v", login, token, err)
					//return nil, fmt.Errorf("cant collect campaings with parameters: collectCampaings(%s, %s), error: %v", login.Login, token, err)
				}
				camps = append(camps, result)
				logrus.Info("Inside gorouitne append is ok for login:", login)
			}
		}()
	}
	for _, c := range resultA.Clients {

		chAC <- c // add agClient to the queue

	}
	close(chAC) // This tells the goroutines there's nothing else to do
	wg.Wait()   // Wait for the threads to finish

	fmt.Println("collectCampaingsfromAgency END", time.Now())
	return camps, nil
}

func collectCampaingsfromAgencyConcurently2(login, token string) (res []campaigns.GetResponse, err error) {
	fmt.Println("collectCampaingsfromAgencyConcurently START", time.Now())
	var YandexConnectionsLimit = 4
	var wg sync.WaitGroup
	wg.Add(YandexConnectionsLimit)
	resultA, err := collectAgencyClients(login, token)
	if err != nil {
		logr.Errorln("collectCampaingsfromAgency  error: ", err)
		return
	}
	var camps []campaigns.GetResponse
	for _, ag := range resultA.Clients {
		go func() {
			for {
				for _, login := range ag.Representatives {
					result, err := collectCampaings(login.Login, token)
					if err != nil {
						logr.Errorln("cant collect campaings with parameters: collectCampaings(%s, %s) error: %v", login.Login, token, err)
						//return nil, fmt.Errorf("cant collect campaings with parameters: collectCampaings(%s, %s), error: %v", login.Login, token, err)
					}
					camps = append(camps, result)
				}
				wg.Done()
			}
		}()
	}
	wg.Wait() // Wait for the threads to finish

	fmt.Println("collectCampaingsfromAgency END", time.Now())
	return camps, nil
}
func collectCampaings(login, token string) (res campaigns.GetResponse, err error) {
	client := api.NewClient()
	//clients.New()"f20.ru", "API_TOKEN"
	client.Token = api.Token{Value: token}
	client.Login = login
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

	service := campaigns.New(client)

	result, err := service.Get(campaingsInfo)
	if err != nil {
		return res, fmt.Errorf("collectCampaings service.Get error %v", err)
	}

	return result, nil
}
func collectClientInfo(login, token string) (res clients.GetResponse, err error) {
	client := api.NewClient()
	//clients.New()"f20.ru", "API_TOKEN"
	client.Token = api.Token{Value: token}
	client.Login = login
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
