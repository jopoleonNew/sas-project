package yandex

import (
	"log"
	"net/http"
	"sync"

	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"
	yad "gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

// GetYandexAccessToken handles requests from Yandex Direct Api to application
// and parsing URL of yandex request to extract confirmation code and accountlogin, then makes
// MakeYandexOauthRequest() to get access token and then saves access token to database.
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

	oauthresp, err := yad.MakeYandexOauthRequest(code)
	if err != nil {
		log.Println("GetYandexAccessToken MakeYandexOauthRequest error: ", err)
		w.Write([]byte("GetYandexAccessToken MakeYandexOauthRequest error: " + err.Error()))
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
		log.Println("..................//////////Hello agency ", oauthresp.AccessToken)
		account := yad.NewAccount()
		account.Login = accinfo.Accountlogin
		account.OAuthToken = oauthresp.AccessToken
		log.Println("..................//////////Hello agency ", oauthresp.AccessToken)
		agencystruct, err := account.GetAgencyLogins()
		if err != nil {
			log.Println("SubmitConfirmationYandexCode GetAgencyLogins error: ", err)
			return
		}
		log.Println("SubmitConfirmationYandexCode GetAccountInfo: ", accinfo)
		log.Println("SubmitConfirmationYandexCode account.GetAgencyLogins(): ", agencystruct)

		//DONETODO: add concurrency to getting campaign list for every agency account login
		wg := sync.WaitGroup{}
		wg.Add(len(agencystruct))
		for _, agClient := range agencystruct {
			go func() {
				agencyacc := model.NewAccount()
				agencyacc.Accountlogin = agClient.Login
				agencyacc.Username = username
				agencyacc.Email = agClient.Representatives[0].Email
				agencyacc.YandexRole = agClient.Representatives[0].Role
				agencyacc.Source = "Яндекс Директ"
				agencyacc.OauthToken = oauthresp.AccessToken
				//var campjson model.CampaingsGetResult
				account := yad.NewAccount()
				account.Login = agClient.Login
				account.OAuthToken = accinfo.OauthToken
				yadcamps, err := account.GetCampaignList()
				if err != nil {
					log.Fatal("SubmitConfirmationYandexCode GetAgencyLogins GetCampaignList: ", err)
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

				err = agencyacc.AdvanceUpdate()
				if err != nil {
					log.Fatal("SubmitConfirmationYandexCode agencyacc.Update() error: ", err)
					return
				}
				wg.Done()
			}()
		}
		wg.Wait()
		return
	}
	yadacc := yad.NewAccount()
	yadacc.Login = accountlogin
	yadacc.OAuthToken = oauthresp.AccessToken
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
	err = acc.AdvanceUpdate()
	if err != nil {
		log.Fatal("SubmitConfirmationYandexCode acc.AdvanceUpdate() error: ", err)
		return
	}
	//log.Println("SubmitConfirmationYandexCode account change info: ", chInfo)
	//campsbyte, err := json.Marshal(yadcamps)
	//if err != nil {
	//	log.Println("SubmitConfirmationYandexCode json.Marshal error: ", err)
	//	return
	//}
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
	return
}
