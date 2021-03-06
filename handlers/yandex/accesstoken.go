package yandex

import (
	"log"
	"net/http"

	"sync"

	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

// GetYandexAccessToken handles requests from Yandex Direct Api to application
// and parsing URL of yandex request to extract confirmation code and accountlogin, then makes
// GetYandexToken() to get access token and then saves access token to database.
// /getyandexaccesstoken endpoint
func GetYandexAccessToken(w http.ResponseWriter, r *http.Request) {
	log.Println("GetYandexAccessToken used")
	query := r.URL.Query()
	log.Println("GetYandexAccessToken income URL query: ", r.URL.Query())

	s := query["code"]
	if s == nil || len(s) == 0 {
		log.Println("Request from YandexOuath received without code")
		w.Write([]byte("Request from YandexOuath received without code"))
		return
	}
	// "state" is the Yandex account login sent with GetAuthCodeYandexHandler()
	al := query["state"]
	if al == nil || len(al) == 0 {
		log.Println("Request from YandexOuath received without state")
		w.Write([]byte("Request from YandexOuath received without state"))
		return
	}
	accountlogin := al[0]
	code := s[0]
	log.Println("Token: ", code, "Accountlogin: ", accountlogin)
	//for better testing
	oauthresp, err := yad.GetYandexToken(code, yad.API_YANDEX_OAUTH_URL)
	if err != nil {
		log.Println("GetYandexAccessToken GetYandexToken error: ", err)
		w.Write([]byte("GetYandexAccessToken GetYandexToken error: " + err.Error()))
		return
	}
	username, err := utils.GetUsernamefromRequestSession(r)
	if err != nil {
		log.Println("SubmitConfirmationYandexCode getUsernamefromRequestSession: ", err)
		return
	}

	// saving access token to DB: oauthresp.AccessToken
	acc := model.NewAccount()
	acc.Accountlogin = accountlogin
	acc.Username = username
	acc.Source = "Яндекс Директ"
	acc.OauthToken = oauthresp.AccessToken
	acc.AdvanceUpdate()
	//log.Println("SubmitConfirmationYandexCode: ", string(data))
	accinfo, err := acc.GetInfo()
	if err != nil {
		log.Println("SubmitConfirmationYandexCode acc.GetInfo() error: ", err)
		w.Write([]byte("SubmitConfirmationYandexCode acc.GetInfo() error: " + err.Error()))
		return
	}
	log.Println("SubmitConfirmationYandexCode GetAccountInfo: ", accinfo)

	if accinfo.YandexRole == "agency" {
		log.Println("..................//////////Hello agency ", acc.Username, oauthresp.AccessToken)
		account := yad.NewAccount(accinfo.Accountlogin, oauthresp.AccessToken)
		agencystruct, err := account.GetAgencyLogins()
		if err != nil {
			log.Println("SubmitConfirmationYandexCode GetAgencyLogins error: ", err)
			return
		}
		YandexConnectionsLimit := 10
		chAC := make(chan yad.ClientAg, 3) // This number 3 can be anything as long as it's larger than YandexConnectionsLimit
		var wg sync.WaitGroup
		// This starts number of goroutines that wait for add new account
		// of agency and its campaings to DB
		wg.Add(YandexConnectionsLimit)
		for i := 0; i < YandexConnectionsLimit; i++ {
			go func() {
				for {
					agClient, ok := <-chAC
					if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
						wg.Done()
						return
					}
					//log15.Info("Inside agency handling for loop ", "agency", agClient)
					log.Println("----------.......///// \n Inside agency handling for loop ", agClient)
					agencyacc := model.NewAccount()
					agencyacc.Accountlogin = agClient.Login
					agencyacc.Username = username
					agencyacc.Email = agClient.Representatives[0].Email
					agencyacc.YandexRole = agClient.Representatives[0].Role
					agencyacc.Source = "Яндекс Директ"
					agencyacc.OauthToken = oauthresp.AccessToken
					//var campjson model.CampaingsGetResult
					account := yad.NewAccount(agClient.Login, accinfo.OauthToken)
					yadcamps, err := account.GetCampaignList()
					if err != nil {
						w.Write([]byte("SubmitConfirmationYandexCode GetAgencyLogins GetCampaignList err:" + err.Error()))
						log.Fatal("SubmitConfirmationYandexCode GetAgencyLogins GetCampaignList err: ", err)
						//w.Write([]byte("SubmitConfirmationYandexCode GetCampaignsListYandex:" + err.Error()))
						return
					}
					acccamps := make([]model.Campaign, len(yadcamps))
					for i, camp := range yadcamps {
						acccamps[i].ID = camp.ID
						acccamps[i].Status = camp.Status
						acccamps[i].Name = camp.Name
					}
					agencyacc.CampaignsInfo = acccamps
					acc.AgencyClients = append(acc.AgencyClients, agClient.Login)
					err = agencyacc.AdvanceUpdate()
					if err != nil {
						log.Fatal("SubmitConfirmationYandexCode agencyacc.Update() error: ", err)
						return
					} // do the thing
				}
			}()
		}

		// Now the jobs can be added to the channel, which is used as a queue
		for _, a := range agencystruct {
			chAC <- a // add agClient to the queue
		}

		close(chAC) // This tells the goroutines there's nothing else to do
		wg.Wait()   // Wait for the threads to finish

		err = acc.AdvanceUpdate()
		if err != nil {
			log.Fatal("SubmitConfirmationYandexCode acc.AdvanceUpdate() error: ", err)
			return
		}
		return
	}
	yadacc := yad.NewAccount(accountlogin, oauthresp.AccessToken)
	yadcamps, err := yadacc.GetCampaignList()
	if err != nil {
		log.Println("SubmitConfirmationYandexCode GetCampaignsListYandex: ", err)
		w.Write([]byte("SubmitConfirmationYandexCode GetCampaignsListYandex:" + err.Error()))
		return
	}

	//log.Println("\n\n Campaings slice: ", yadcamps)
	acccamps := make([]model.Campaign, len(yadcamps))
	for i, camp := range yadcamps {
		acccamps[i].ID = camp.ID
		acccamps[i].Status = camp.Status
		acccamps[i].Name = camp.Name
	}
	acc.CampaignsInfo = acccamps
	//updating acc status on active
	acc.Status = "active"

	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
	return
}
