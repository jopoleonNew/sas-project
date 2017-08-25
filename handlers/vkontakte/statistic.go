package vkontakte

import (
	"fmt"
	"html/template"
	"net/http"

	"log"

	"strconv"

	"strings"

	"time"

	"encoding/json"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
	vk "gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
)

const ctLayout = "2006-01-02"

type ss struct {
	Id          int
	Day         string
	Spent       string
	Impressions int
	Clicks      int
	Reach       int
}

func CollectVKStatistic(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	logrus.Println("CollectVKStatistic used with username: ", username)

	//if request method is GET - parsing html template and sending it to response
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
	//if request method is GET - getting statistic from third-party API
	if r.Method == "POST" {
		logrus.Info("CollectVKStatistic  POST ", username)
		startdate := r.FormValue("startdate")
		enddate := r.FormValue("enddate")
		sttime, err := time.Parse(ctLayout, startdate)
		if err != nil {
			logrus.Error("CollectVKStatistic time.Parse error: ", err)
			http.Error(w, fmt.Sprintf("cant parse recieved time value, :%v", err), http.StatusConflict)
			return
		}
		endtime, err := time.Parse(ctLayout, enddate)
		if err != nil {
			logrus.Error("CollectVKStatistic time.Parse error: ", err)
			http.Error(w, fmt.Sprintf("cant parse recieved time value, :%v", err), http.StatusConflict)
			return
		}
		fomrmatstart := sttime.Format(ctLayout)
		formatend := endtime.Format(ctLayout)
		log.Println(" Recieved time values: ", startdate, enddate)
		query := r.URL.Query()
		if query["login"] == nil || len(query["login"]) == 0 {
			logrus.Error("CollectVKStatistic GET request recieved without acclount login. ")
			http.Error(w, fmt.Sprintf("GET request recieved without acclount login. "), http.StatusBadRequest)
			return
		}
		accountlogin := query["login"][0]
		logrus.Infof("CollectVKStatistic used with POST: ", username, accountlogin)
		a := model.NewAccount2("", "", accountlogin, "")
		info, err := a.GetInfo()
		if err != nil {
			logrus.Errorf("ReportTemplateHandler a.GetInfo() %v error: %v", a, err)
			http.Error(w, fmt.Sprintf("can't find in db Yandex account %v \n error: %+v:", a, err), http.StatusBadRequest)
			return
		}
		var adids []string
		for _, c := range info.CampaignsInfo {
			for _, ad := range c.Ads {
				adids = append(adids, strconv.Itoa(ad.ID))
			}

		}
		if len(adids) < 1999 {
			p := make(map[string]string)
			p["account_id"] = a.Accountlogin
			p["ids_type"] = "ad"
			p["ids"] = strings.Join(adids, ", ")
			//p["ids_type"] = "campaign"
			p["period"] = "day"
			p["date_from"] = fomrmatstart
			p["date_to"] = formatend
			res, err := collectStatistic(info.AuthToken, p)
			if err != nil {
				logrus.Errorf("CollectVKStatistic collectStatistic error: %v", err)
				w.Write([]byte("CollectVKStatistic collectStatistic error: " + err.Error()))
				return
			}
			//logrus.Infof("Result from Statistic: %s", res)
			w.Write(res)
		} else {
			logrus.Errorf("Amount of campings IDs is more than 2000: %d", len(adids))
		}
	}

}
func collectStatistic(token string, params map[string]string) (vk.AdStatistic, error) {
	var stats vk.AdStatistic
	if token == "" {
		return stats, fmt.Errorf("Token is empty")
	}
	//var camps vk.AdsCampaigns
	resp, err := vk.Request(token, "ads.getStatistics", params)
	if err != nil {
		logrus.Errorf("collectStatistic vk.Request error: %v", err)
		return stats, fmt.Errorf("collectStatistic vk.Request error: %v", err)
	}

	if err := json.Unmarshal(resp, &stats); err != nil {
		logrus.Errorf("can't unmarshal VK response from ads.getAds, error: &v", err)

		return stats, fmt.Errorf("collectAds json.Unmarshal error: %v", err)
	}

	return stats, nil
}
