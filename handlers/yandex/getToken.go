package yandex

import (
	//"log"
	"net/http"

	"fmt"

	"strings"

	"time"

	"sync"

	"os"

	"encoding/json"

	api "github.com/nk2ge5k/goyad"
	"github.com/nk2ge5k/goyad/agencyclients"
	"github.com/nk2ge5k/goyad/campaigns"
	"github.com/nk2ge5k/goyad/clients"
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
	logr.Info(" GetYandexToken Token: ", code, " Accountlogin: ", accountlogin)

	oauthresp, err := yad.MakeYandexOauthRequest(code)
	if err != nil {
		logr.Println("GetYandexAccessToken MakeYandexOauthRequest error: ", err)
		w.Write([]byte("GetYandexAccessToken MakeYandexOauthRequest error: " + err.Error()))
		return
	}
	result, err := collectCampaings(accountlogin, oauthresp.AccessToken)
	if err != nil {
		if strings.Contains(err.Error(), "53") {
			result, err := collectCampaingsfromAgencyConcurently(accountlogin, oauthresp.AccessToken)
			if err != nil {
				logr.Println("GetYandexAccessToken collectCampaingsfromAgency(%s, %s) error: ", accountlogin, oauthresp.AccessToken, err)
				http.Error(w, fmt.Sprintf("cant recieve campaings from Yandex.Direct with parametrs %s %s", accountlogin, err), http.StatusBadRequest)
				return
			}
			logrus.Info("collectCampaingsfromAgency SUCCESS")
			//f, err := os.Create("AgencyData.json")
			//if err != nil {
			//	panic(err)
			//}
			//defer f.Close()
			//agbytes, err := json.Marshal(result)
			//if err != nil {
			//	panic(err)
			//}
			//n2, err := f.Write(agbytes)
			//if err != nil {
			//	panic(err)
			//}

			fmt.Printf("wrote %d bytes\n", n2)

		} else {
			logr.Errorln("GetYandexAccessToken collectCampaings error: ", err)
			http.Error(w, fmt.Sprintf("cant recieve campaings from Yandex.Direct %s", err), http.StatusBadRequest)
			return
		}

	}

	//logrus.Printf("Result from Yandex about %s: \n %+v", accountlogin, result.Campaigns[5])

	acccamps := model.AdaptYandexCampaings(result)

	for _, c := range acccamps {
		logrus.Printf("Result from Yandex after adaption: \n %+v", c)
	}
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)

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

	var YandexConnectionsLimit = 10
	chAC := make(chan string, 4) // This number 3 can be anything as long as it's larger than YandexConnectionsLimit
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
				result, err := collectCampaings(login, token)
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

		chAC <- c.Login // add agClient to the queue

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
			"AccountQuality", "Archived", "ClientId", "ClientInfo", "CountryId", "CreatedAt", "Currency", "Grants", "Login", "Notification", "OverdraftSumAvailable", "Phone", "Representatives", "Restrictions", "Settings", "Type", "VatRate",
		},
	}
	service2 := clients.New(&client)
	result, err := service2.Get(clientInfo)
	if err != nil {
		return res, fmt.Errorf("collectClientInfo service.Get error %v", err)
	}
	return result, nil
}

//YandexConnectionsLimit := 10
//chAC := make(chan yad.Client, 3) // This number 3 can be anything as long as it's larger than YandexConnectionsLimit
//var wg sync.WaitGroup
//
//// This starts number of goroutines that wait for add new account
//// of agency and its campaings to DB
//wg.Add(YandexConnectionsLimit)
//for i := 0; i < YandexConnectionsLimit; i++ {
//go func() {
//for {
//agClient, ok := <-chAC
//if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
//wg.Done()
//return
//}
////log15.Info("Inside agency handling for loop ", "agency", agClient)
//log.Println("------------............. ///// \n Inside agency handling for loop ", agClient)
//agencyacc := model.NewAccount()
//agencyacc.Accountlogin = agClient.Login
//agencyacc.Username = username
//agencyacc.Email = agClient.Representatives[0].Email
//agencyacc.YandexRole = agClient.Representatives[0].Role
//agencyacc.Source = "Яндекс Директ"
//agencyacc.OauthToken = oauthresp.AccessToken
////var campjson model.CampaingsGetResult
//account := yad.NewAccount(agClient.Login, accinfo.OauthToken)
//yadcamps, err := account.GetCampaignList()
//if err != nil {
//w.Write([]byte("SubmitConfirmationYandexCode GetAgencyLogins GetCampaignList err:" + err.Error()))
//log.Fatal("SubmitConfirmationYandexCode GetAgencyLogins GetCampaignList err: ", err)
////w.Write([]byte("SubmitConfirmationYandexCode GetCampaignsListYandex:" + err.Error()))
//return
//}
//acccamps := make([]model.Campaign, len(yadcamps))
//for i, camp := range yadcamps {
//acccamps[i].ID = camp.ID
//acccamps[i].Status = camp.Status
//acccamps[i].Name = camp.Name
//}
//agencyacc.CampaignsInfo = acccamps
//acc.AgencyClients = append(acc.AgencyClients, agClient.Login)
//err = agencyacc.AdvanceUpdate()
//if err != nil {
//log.Fatal("SubmitConfirmationYandexCode agencyacc.Update() error: ", err)
//return
//} // do the thing
//}
//}()
//}
//
//// Now the jobs can be added to the channel, which is used as a queue
//for _, a := range agencystruct {
//chAC <- a // add agClient to the queue
//}
//
//close(chAC) // This tells the goroutines there's nothing else to do
//wg.Wait()   // Wait for the threads to finish
