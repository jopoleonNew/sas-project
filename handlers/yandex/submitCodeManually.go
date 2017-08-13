package yandex

import (
	"encoding/json"
	"log"
	"net/http"

	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

// SubmitConfirmationYandexCode handles requests from client to submiting confirmation Yandex Code
// then parses client's accountlogin in Yandex and uses GetYandexToken to get access token
// @impotrant This method is used for manual receiving of confirmation Yandex Code
func SubmitConfirmationYandexCode(w http.ResponseWriter, r *http.Request) {

	_, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println("GetStatSliceHandler store.Get error:", err)
		w.Write([]byte("GetStatSliceHandler store.Get error " + err.Error()))
		return
	}
	log.Println("SubmitConfirmationYandexCode used")

	r.ParseForm()
	yandexcode := r.FormValue("yandexcode")
	accountlogin := r.FormValue("accountlogin")

	//log.Println("SubmitConfirmationYandexCode: ", yandexcode, accountlogin)

	oauthresp, err := yad.GetYandexToken(yandexcode, yad.API_YANDEX_OAUTH_URL)
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

	// seting accsess token to DB
	acc := model.NewAccount()
	acc.Accountlogin = accountlogin
	acc.Username = username
	acc.Source = "Яндекс Директ"
	acc.OauthToken = oauthresp.AccessToken

	//log.Println("..................//////////Hello agency ", oauthresp.AccessToken)

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

		// if account who submiting code is agency,
		// we must obtain all agency clients info(their login, email, etc.)
		account := yad.NewAccount(accountlogin, oauthresp.AccessToken)
		agencystruct, err := account.GetAgencyLogins()
		if err != nil {
			log.Println("SubmitConfirmationYandexCode GetAgencyLogins error: ", err)
			return
		}
		// adding client's logins of agency to DB
		agencylogins := make([]string, len(agencystruct))
		for i, client := range agencystruct {
			agencylogins[i] = client.Login
		}
		acc.AgencyClients = agencylogins
		acc.AdvanceUpdate()

		log.Println("SubmitConfirmationYandexCode GetAccountInfo: ", accinfo)
		log.Println("SubmitConfirmationYandexCode account.GetAgencyLogins(): ", agencystruct)
		//acc.AgencyClients = agencystruct
		//DONETODO: add concurrency to getting campaign list for every agency account login

		// inside this loop we get campaigns for all agency clients
		// and create for every of them account in DB with
		for _, agClient := range agencystruct {

			agencyacc := model.NewAccount()
			agencyacc.Accountlogin = agClient.Login
			agencyacc.Username = username
			agencyacc.Email = agClient.Representatives[0].Email
			agencyacc.YandexRole = agClient.Representatives[0].Role
			agencyacc.Source = "Яндекс Директ"
			agencyacc.OauthToken = oauthresp.AccessToken
			//var campjson model.CampaingsGetResult
			account := yad.NewAccount(agClient.Login, oauthresp.AccessToken)
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
	acccamps := make([]model.Campaign, len(yadcamps))
	for i, camp := range yadcamps {
		acccamps[i].ID = camp.ID
		acccamps[i].Status = camp.Status
		acccamps[i].Name = camp.Name
	}
	acc.CampaignsInfo = acccamps
	//updating account status on active
	acc.Status = "active"
	err = acc.AdvanceUpdate()
	if err != nil {
		log.Fatal("SubmitConfirmationYandexCode acc.AdvanceUpdate() error: ", err)
		return
	}
	//log.Println("SubmitConfirmationYandexCode account change info: ", chInfo)
	campsbyte, err := json.Marshal(yadcamps)
	if err != nil {
		log.Println("SubmitConfirmationYandexCode json.Marshal error: ", err)
		return
	}
	w.Write(campsbyte)
}
