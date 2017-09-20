package yandex

import (
	"encoding/json"
	"log"
	//"model"
	"net/http"
	"strconv"

	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"
	yad "gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
	//"utils"
)

// GetCampaingStatsHandler handling requests from client to get statistic about
// sent campaing ID's sending request to Yandex Direct API
func GetCampaingStatsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetCampaingStatsHandler used")
	r.ParseForm()
	username := r.FormValue("username")
	accountlogin := r.FormValue("accountlogin")
	campaingId := r.FormValue("campaingId")
	log.Println("GetCampaingStatsHandler recieved info: ", username, accountlogin, campaingId)
	acc := model.NewAccount()
	acc.Username = username
	acc.Accountlogin = accountlogin
	accinfo, err := acc.GetInfo()
	if err != nil {
		log.Println("GetCampaingStatsHandler acc.GetInfo() error: ", err)
		w.Write([]byte("GetCampaingStatsHandler acc.GetInfo() error: " + err.Error()))
		return
	}
	log.Println("accinfo: ", accinfo)
	//from Yandex Direct API
	//StartDate	Начальная дата отчетного периода, за который возвращается статистика (YYYY-MM-DD).
	//EndDate	Конечная дата отчетного периода, за который возвращается статистика (YYYY-MM-DD).
	startdate := "2017-04-04"
	enddate := "2017-04-20"

	idint, err := strconv.Atoi(campaingId)
	if err != nil {
		log.Println("GetCampaingStatsHandler strconv.Atoi error: ", err)
		w.Write([]byte("GetCampaingStatsHandler strconv.Atoi error: " + err.Error()))
		return
	}
	account := yad.NewAccount(accinfo.Accountlogin, accinfo.OauthToken)
	statres, err := account.GetStatistics([]int{idint}, startdate, enddate)
	if err != nil {
		log.Println("GetCampaingStatsHandler GetStatistics error: ", err)
		w.Write([]byte("GetCampaingStatsHandler GetStatistics error: " + err.Error()))
		return
	}
	statbytes, err := json.Marshal(statres)
	if err != nil {
		log.Println("GetCampaingStatsHandler json.Marshal error: ", err)
		w.Write([]byte("GetCampaingStatsHandler json.Marshal error: " + err.Error()))
		return
	}
	w.Write(statbytes)
}

// RefreshCampaignsListHandler handles client's requests to
// refresh campaigns in DB with given username and account login in FormValue
func RefreshCampaignsListHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("RefreshCampaignsListHandler used")
	username, err := utils.GetUsernamefromRequestSession(r)
	if err != nil {
		log.Println("RefreshCampaignsListHandler getUsernamefromRequestSession error: ", err)
		w.Write([]byte("RefreshCampaignsListHandler getUsernamefromRequestSession error: " + err.Error()))
		return
	}
	r.ParseForm()
	accountlogin := r.FormValue("accountlogin")

	acc := model.NewAccount()
	acc.Username = username
	acc.Accountlogin = accountlogin
	acc.Source = "Яндекс Директ"
	accinfo, err := acc.GetInfo()
	if err != nil {
		log.Println("RefreshCampaignsListHandler acc.GetInfo() error: ", err)
		w.Write([]byte("RefreshCampaignsListHandler acc.GetInfo() error: " + err.Error()))
		return
	}
	yadacc := yad.NewAccount(accinfo.Accountlogin, accinfo.OauthToken)
	yadcamps, err := yadacc.GetCampaignList()
	if err != nil {
		log.Println("RefreshCampaignsListHandler GetCampaignList error: ", err)
		w.Write([]byte("RefreshCampaignsListHandler GetCampaignList error: " + err.Error()))
		return
	}
	//log.Println("RefreshCampaignsListHandler campaings: ", yadcamps)

	//reassigning yandex campaigns to model application campaigns struct
	acccamps := make([]model.Campaign, len(yadcamps))
	for i, camp := range yadcamps {
		acccamps[i].ID = camp.ID
		acccamps[i].Status = camp.Status
		acccamps[i].Name = camp.Name
	}
	acc.CampaignsInfo = acccamps
	err = acc.AdvanceUpdate()
	if err != nil {
		log.Println("RefreshCampaignsListHandler acc.Update() error: ", err)
		w.Write([]byte("RefreshCampaignsListHandler acc.Update() error: " + err.Error()))
		return
	}
	//log.Println("RefreshCampaignsListHandler account change info: ", chInfo)
	w.Write([]byte("Аккаунты в базе данных обновленны в соответствии с Яндекс.Директом"))
}
