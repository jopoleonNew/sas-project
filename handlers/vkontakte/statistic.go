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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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
			logrus.Error(err)
			fmt.Fprintf(w, err.Error())
		}

		err = t.Execute(w, data)
		if err != nil {
			logrus.Error(err)
			fmt.Fprintf(w, err.Error())
		}
	}
	//if request method is POST - getting statistic from db "statistic"
	if r.Method == "POST" {
		logrus.Info("CollectVKStatistic  POST ", username)
		query := r.URL.Query()
		if query["login"] == nil || len(query["login"]) == 0 {
			logrus.Error("CollectVKStatistic GET request recieved without acclount login. ")
			http.Error(w, fmt.Sprintf("GET request recieved without acclount login. "), http.StatusBadRequest)
			return
		}
		accountlogin := query["login"][0]
		res, err := model.GetVKStatistic(accountlogin)
		accinfo, err := model.NewAccount2("", "Vkontakte", accountlogin, "").GetInfo()
		if err != nil {
			logrus.Errorf("CollectVKStatistic GetInfo() error: %v", err)
			w.Write([]byte("CollectVKStatistic GetInfo() error: " + err.Error()))
			return
		}
		//forming struct to return to user on his request
		giveaway := []struct {
			ID               int
			Name             string
			Type             string
			Day              string      `json:"day"`
			Spent            string      `json:"spent,omitempty"`
			Impressions      interface{} `json:"impressions,omitempty"`
			Clicks           int         `json:"clicks,omitempty"`
			VideoViews       int         `json:"video_views,omitempty"`
			VideoViews_half  int         `json:"video_views_half,omitempty"`
			VideoViews_full  int         `json:"video_views_full,omitempty"`
			VideoClicks_site int         `json:"video_clicks_site,omitempty"`
			JoinRate         int         `json:"join_rate,omitempty"`
		}{}
		for _, stat := range res.Response {
			for _, c := range accinfo.CampaignsInfo {
				for _, ad := range c.Ads {
					if ad.ID == stat.ID {
						for _, s := range stat.Stats {
							giveaway = append(giveaway, struct {
								ID               int
								Name             string
								Type             string
								Day              string      `json:"day"`
								Spent            string      `json:"spent,omitempty"`
								Impressions      interface{} `json:"impressions,omitempty"`
								Clicks           int         `json:"clicks,omitempty"`
								VideoViews       int         `json:"video_views,omitempty"`
								VideoViews_half  int         `json:"video_views_half,omitempty"`
								VideoViews_full  int         `json:"video_views_full,omitempty"`
								VideoClicks_site int         `json:"video_clicks_site,omitempty"`
								JoinRate         int         `json:"join_rate,omitempty"`
							}{
								ID:               stat.ID,
								Name:             ad.Name,
								Type:             stat.Type,
								Day:              s.Day,
								Spent:            s.Spent,
								Impressions:      s.Impressions,
								Clicks:           s.Clicks,
								VideoViews:       s.VideoViews,
								VideoViews_half:  s.VideoViews_half,
								VideoViews_full:  s.VideoViews_full,
								VideoClicks_site: s.VideoClicks_site,
								JoinRate:         s.JoinRate,
							})
						}

					}
				}
			}
		}
		//logrus.Info(giveaway)
		bres, err := json.Marshal(giveaway)
		if err != nil {
			logrus.Errorf("CollectVKStatistic json.Marshal(res) error: %v", err)
			w.Write([]byte("CollectVKStatistic json.Marshal(res) error: " + err.Error()))
			return
		}
		w.Write(bres)
		//} else {
		//	logrus.Errorf("Amount of campings IDs is more than 2000: %d", len(adids))
		//	w.Write([]byte("Amount of campings IDs is more than 2000"))
		//	return
		//}
	}

}

func collectStatistic(token string, params map[string]string) (vk.AdStatistic, error) {
	var stats vk.AdStatistic
	if token == "" {
		return stats, fmt.Errorf("token is empty")
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

func CollectVKStatistic_old(w http.ResponseWriter, r *http.Request) {
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
			logrus.Error(err)
			fmt.Fprintf(w, err.Error())
		}

		err = t.Execute(w, data)
		if err != nil {
			logrus.Error(err)
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
		logrus.Infof("CollectVKStatistic used with POST: %v %v", username, accountlogin)
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
			//res, err := collectStatistic(info.AuthToken, p)
			//if err != nil {
			//	logrus.Errorf("CollectVKStatistic collectStatistic error: %v", err)
			//	w.Write([]byte("CollectVKStatistic collectStatistic error: " + err.Error()))
			//	return
			//}
			//logrus.Infof("Result from Statistic: %s", res)
			//err = model.SaveVKStatistic(a.Accountlogin, res)
			//if err != nil {
			//	logrus.Errorf("CollectVKStatistic model.SaveVKStatistic error: %v", err)
			//	w.Write([]byte("CollectVKStatistic model.SaveVKStatistic error: " + err.Error()))
			//	return
			//}

			res, err := model.GetVKStatistic(a.Accountlogin)
			if err != nil {
				logrus.Errorf("CollectVKStatistic GetVKStatistic error: %v", err)
				w.Write([]byte("CollectVKStatistic GetVKStatistic error: " + err.Error()))
				return
			}
			accinfo, err := model.NewAccount2("", "vkontakte", "a.Accountlogin", "").GetInfo()
			if err != nil {
				logrus.Errorf("CollectVKStatistic GetInfo() error: %v", err)
				w.Write([]byte("CollectVKStatistic GetInfo() error: " + err.Error()))
				return
			}
			giveaway := []struct {
				ID               int
				Name             string
				Type             string
				Day              string      `json:"day"`
				Spent            string      `json:"spent,omitempty"`
				Impressions      interface{} `json:"impressions,omitempty"`
				Clicks           int         `json:"clicks,omitempty"`
				VideoViews       int         `json:"video_views,omitempty"`
				VideoViews_half  int         `json:"video_views_half,omitempty"`
				VideoViews_full  int         `json:"video_views_full,omitempty"`
				VideoClicks_site int         `json:"video_clicks_site,omitempty"`
				JoinRate         int         `json:"join_rate,omitempty"`
			}{}
			for _, stat := range res.Response {
				for _, c := range accinfo.CampaignsInfo {
					for _, ad := range c.Ads {
						if ad.ID == stat.ID {
							for _, s := range stat.Stats {
								giveaway = append(giveaway, struct {
									ID               int
									Name             string
									Type             string
									Day              string      `json:"day"`
									Spent            string      `json:"spent,omitempty"`
									Impressions      interface{} `json:"impressions,omitempty"`
									Clicks           int         `json:"clicks,omitempty"`
									VideoViews       int         `json:"video_views,omitempty"`
									VideoViews_half  int         `json:"video_views_half,omitempty"`
									VideoViews_full  int         `json:"video_views_full,omitempty"`
									VideoClicks_site int         `json:"video_clicks_site,omitempty"`
									JoinRate         int         `json:"join_rate,omitempty"`
								}{
									ID:               stat.ID,
									Name:             ad.Name,
									Type:             stat.Type,
									Day:              s.Day,
									Spent:            s.Spent,
									Impressions:      s.Impressions,
									Clicks:           s.Clicks,
									VideoViews:       s.VideoViews,
									VideoViews_half:  s.VideoViews_half,
									VideoViews_full:  s.VideoViews_full,
									VideoClicks_site: s.VideoClicks_site,
									JoinRate:         s.JoinRate,
								})
							}

						}
					}
				}
			}
			bres, err := json.Marshal(giveaway)
			if err != nil {
				logrus.Errorf("CollectVKStatistic json.Marshal(res) error: %v", err)
				w.Write([]byte("CollectVKStatistic json.Marshal(res) error: " + err.Error()))
				return
			}
			w.Write(bres)
		} else {
			logrus.Errorf("Amount of campings IDs is more than 2000: %d", len(adids))
			w.Write([]byte("Amount of campings IDs is more than 2000"))
			return
		}
	}

}
