package yandex

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

const ctLayout = "2006-01-02"

func GetStatSliceHandler(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("username").(string)

	log.Println("GetStatSliceHandler used by", username)
	startdate := r.FormValue("startdate")
	enddate := r.FormValue("enddate")
	//dataLayout :="2006-01-02"
	//log.Println("GetStatSliceHandler ", username)
	sttime, err := time.Parse(ctLayout, startdate)
	if err != nil {
		log.Fatal("time.Parse error: ", err)
		return
	}
	endtime, err := time.Parse(ctLayout, enddate)
	if err != nil {
		log.Fatal("time.Parse error: ", err)
		return
	}
	var data model.TemplateInfoStruct
	data.CurrentUser = username
	//acc := model.NewAccount()
	//acc.Username = username
	user := model.NewUser()
	user.Username = username

	//acclist, err := acc.GetInfoList()
	acclist, err := user.GetAccountList()
	if err != nil {
		log.Println("GetStatSliceHandler acc.GetInfoList() error:", err)
		w.Write([]byte("GetStatSliceHandler error: " + err.Error()))
		return
	}
	//log.Println("GetStatSliceHandlerGetStatSliceHandler acclist ", acclist)
	var statsslice yad.CampaignStatSlice

	//DONETODO: Add concurrency for every account in loop
	//wg := sync.WaitGroup{}
	for _, camp := range acclist {

		if camp.Source == "Яндекс Директ" {

			var idslice []int
			if camp.YandexRole == "agency" {
				acc := model.NewAccount()
				acc.Username = username
				acc.Accountlogin = camp.Accountlogin
				agencyInfo, err := acc.GetInfo()
				if err != nil {
					log.Println("GetStatSliceHandler agencyInfo.GetInfo error:", err)
					w.Write([]byte("GetStatSliceHandler agencyInfo.GetInfo error: " + err.Error()))
					return
				}
				for _, agencyAccountLogin := range agencyInfo.AgencyClients {
					agencyAcc := model.NewAccount()
					agencyAcc.Username = username
					agencyAcc.Accountlogin = agencyAccountLogin
					agencyAccInfo, err := agencyAcc.GetInfo()
					if err != nil {
						log.Println("GetStatSliceHandler agencyAccInfo.GetInfo error:", err)
						w.Write([]byte("GetStatSliceHandler agencyAccInfo.GetInfo error: " + err.Error()))
						return
					}
					//log.Println("Inside AccountHandler. Agency's AccountInfo : %+v", agencyAccInfo)
					for _, id := range agencyAccInfo.CampaignsInfo {
						idslice = append(idslice, id.ID)
					}

					//idslice = append(idslice, id.ID)
				}
			}
			for _, id := range camp.CampaignsInfo {
				idslice = append(idslice, id.ID)
			}

			//wg.Add(1)
			//go func() {
			//wg.Wait()
			account := yad.NewAccount(camp.Accountlogin, camp.OauthToken)
			statres, err := account.GetStatisticsConc(idslice, sttime, endtime)
			if err != nil {
				log.Println("GetStatSliceHandlerGetStatSliceHandler GetCampaingsSliceStatistic", err)
				fmt.Fprintf(w, "GetStatSliceHandlerGetStatSliceHandler GetCampaingsSliceStatistic"+err.Error())
				return
			}
			statsslice = append(statsslice, statres...)
			//	wg.Done()
			//}()

		}
	}
	//wg.Wait()
	sort.Sort(statsslice)
	//log.Println("GetStatSliceHandler result stat: ", statsslice)
	reqbytes, err := json.Marshal(statsslice)
	if err != nil {
		log.Println("GetStatSliceHandler json.Marshal error: ", err)
		return
	}
	//log.Println("GetStatSliceHandler AFTER for loop statsslice reqbytes: ", string(reqbytes))
	w.Write(reqbytes)
}

func CollectYandexStatistic(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	log.Println("GetAccountStatistic used by", username)
	if r.Method == "GET" {
		logrus.Info("CollectVKStatistic  GET ", username)
		//vars := mux.Vars(r)
		username := r.Context().Value("username").(string)
		query := r.URL.Query()
		if query["login"] == nil || len(query["login"]) == 0 {
			logrus.Error("CollectVKStatistic GET request recieved without acclount login. ")
			http.Error(w, fmt.Sprintf("GET request recieved without acclount login. "), http.StatusBadRequest)
			return
		}
		accountlogin := query["login"][0]
		logrus.Infof("CollectVKStatistic used with username: %s, accountlogin: %s", username, accountlogin)
		var data model.TemplateInfo
		a := model.NewAccount2("", "", "", "")
		a.Accountlogin = accountlogin
		info, err := a.GetInfo()
		if err != nil {
			logrus.Errorf("ReportTemplateHandler a.GetInfo() %v error: %v", a, err)
			http.Error(w, fmt.Sprintf("can't find in db Yandex account %v \n error: %+v:", a, err), http.StatusBadRequest)
			return
		}
		data.CurrentAccount = info
		data.CurrentUser = username
		t, err := template.New("reports.tmpl").ParseFiles(
			"static/templates/header.tmpl",
			"static/templates/reports.tmpl")
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, err.Error())
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, err.Error())
		}
	}
	if r.Method == "POST" {
		startdate := r.FormValue("startdate")
		enddate := r.FormValue("enddate")
		sttime, err := time.Parse(ctLayout, startdate)
		if err != nil {
			logrus.Error("CollectYandexStatistic time.Parse error: ", err)
			http.Error(w, fmt.Sprintf("cant parse recieved time value, :%v", err), http.StatusConflict)
			return
		}
		endtime, err := time.Parse(ctLayout, enddate)
		if err != nil {
			logrus.Error("CollectYandexStatistic time.Parse error: ", err)
			http.Error(w, fmt.Sprintf("cant parse recieved time value, :%v", err), http.StatusConflict)
			return
		}
		var data model.TemplateInfoStruct
		data.CurrentUser = username
		query := r.URL.Query()
		if query["login"] == nil || len(query["login"]) == 0 {
			logrus.Error("CollectYandexStatistic GET request recieved without acclount login. ")
			http.Error(w, fmt.Sprintf("GET request recieved without acclount login. "), http.StatusBadRequest)
			return
		}
		accountlogin := query["login"][0]
		a := model.NewAccount2("", "", "", "")
		a.Accountlogin = accountlogin
		info, err := a.GetInfo()
		if err != nil {
			logrus.Errorf("CollectYandexStatistic a.GetInfo() %v error: %v", a, err)
			http.Error(w, fmt.Sprintf("can't find in db Yandex account %v \n error: %+v:", a, err), http.StatusBadRequest)
			return
		}
		var idslice []int
		for _, id := range info.CampaignsInfo {
			idslice = append(idslice, id.ID)
		}
		account := yad.NewAccount(info.Accountlogin, info.AuthToken)
		statres, err := account.GetStatisticsConc(idslice, sttime, endtime)
		if err != nil {
			logrus.Errorf("GetYandexAccStat account.GetStatisticsConc %v error: %v", account, err)
			http.Error(w, fmt.Sprintf("can't get statistic from Yandex account %v \n error: %+v:", account, err), http.StatusBadRequest)
			return
		}
		reqbytes, err := json.Marshal(statres)
		if err != nil {
			logrus.Errorf("GetYandexAccStat json.Marshal %v error: %v", statres, err)
			http.Error(w, fmt.Sprintf("can't json.Marshal response from Yandex statistic \n error: %+v:", err), http.StatusBadRequest)
			return
		}
		w.Write(reqbytes)
	}

}
