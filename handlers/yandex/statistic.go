package yandex

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

const ctLayout = "2006-01-02"

func GetStatSliceHandler(w http.ResponseWriter, r *http.Request) {
	_, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println("GetStatSliceHandler store.Get error:", err)
		w.Write([]byte("GetStatSliceHandler store.Get error " + err.Error()))
		return
	}
	///log.Println("Session values Report: ", session.Values)
	username, err := utils.GetUsernamefromRequestSession(r)
	if err != nil {
		log.Println("GetStatSliceHandler GetUsernamefromRequestSession error: ", err)
		w.Write([]byte("GetStatSliceHandler GetUsernamefromRequestSession error: " + err.Error()))
		return
	}
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
