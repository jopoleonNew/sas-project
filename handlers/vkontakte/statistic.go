package vkontakte

import (
	"fmt"
	"html/template"
	"net/http"

	"log"

	"strconv"

	"strings"

	"time"

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
		a := model.NewAccount2("", "", "", "")
		a.Accountlogin = accountlogin
		info, err := a.GetInfo()
		if err != nil {
			logrus.Errorf("ReportTemplateHandler a.GetInfo() %v error: %v", a, err)
			http.Error(w, fmt.Sprintf("can't find in db Yandex account %v \n error: %+v:", a, err), http.StatusBadRequest)
			return
		}
		var ids []string
		for _, c := range info.CampaignsInfo {
			ids = append(ids, strconv.Itoa(c.ID))
		}
		if len(ids) < 1900 {
			p := make(map[string]string)
			p["account_id"] = a.Accountlogin
			p["ids"] = strings.Join(ids, ",")
			p["ids_type"] = "campaign"
			p["period"] = "day"
			p["date_from"] = fomrmatstart
			p["date_to"] = formatend
			res, err := collectStatistic(info.AuthToken, p)
			if err != nil {
				logrus.Errorf("CollectVKStatistic collectStatistic error: %v", err)
				w.Write([]byte("CollectVKStatistic collectStatistic error: " + err.Error()))
				return
			}
			logrus.Infof("Result from Statistic: %s", res)
			w.Write(res)
		} else {
			logrus.Errorf("Amount of campings IDs is more than 1900: %d", len(ids))
		}
	}

}

//account_idидентификатор рекламного кабинета.
//обязательный параметр, целое число
//ids_typeТип запрашиваемых объектов, которые перечислены в параметре ids:
//ad — объявления;
//campaign — кампании;
//client — клиенты;
//office — кабинет.
//обязательный параметр, строка
//idsПеречисленные через запятую id запрашиваемых объявлений, кампаний, клиентов или кабинета, в зависимости от того, что указано в параметре ids_type. Максимум 2000 объектов.
//обязательный параметр, строка
//period Способ группировки данных по датам:
//day — статистика по дням;
//month — статистика по месяцам;
//overall — статистика за всё время;
//Временные ограничения задаются параметрами date_from и date_to.
//обязательный параметр, строка
//date_from Начальная дата выводимой статистики. Используется разный формат дат для разного значения параметра period:
//day: YYYY-MM-DD, пример: 2011-09-27 - 27 сентября 2011
//0 — день создания;
//month: YYYY-MM, пример: 2011-09 - сентябрь 2011
//0 — месяц создания;
//overall: 0
//обязательный параметр, строка
//date_to Конечная дата выводимой статистики. Используется разный формат дат для разного значения параметра period:
//day: YYYY-MM-DD, пример: 2011-09-27 - 27 сентября 2011
//0 — текущий день;
//month: YYYY-MM, пример: 2011-09 - сентябрь 2011
//0 — текущий месяц;
//overall: 0
//обязательный параметр, строка
//Результат
func collectStatistic(token string, params map[string]string) ([]byte, error) {
	if token == "" {
		return []byte{}, fmt.Errorf("Token is empty")
	}
	//var camps vk.AdsCampaigns
	resp, err := vk.Request(token, "ads.getStatistics", params)
	if err != nil {
		logrus.Errorf("collectStatistic vk.Request error: %v", err)
		return resp, fmt.Errorf("collectStatistic vk.Request error: %v", err)
	}
	//logrus.Errorf("VK response from ads.getCampaigns, error: &+v", string(resp))
	//if err := json.Unmarshal(resp, &camps); err != nil {
	//	logrus.Errorf("can't unmarshal VK response from ads.getStatistics, error: &v", err)
	//
	//	return camps, fmt.Errorf("collectStatistic json.Unmarshal error: %v", err)
	//}
	return resp, nil
}
