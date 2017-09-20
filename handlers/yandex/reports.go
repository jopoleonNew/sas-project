package yandex

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"
)

func ReportTemplateHandler(w http.ResponseWriter, r *http.Request) {
	username, err := utils.GetUsernamefromRequestSession(r)
	log.Println("ReportTemplateHandler used by ", username)
	if err != nil {
		log.Println("ReportTemplateHandler GetUsernamefromRequestSession error: ", err)
		w.Write([]byte("ReportTemplateHandler GetUsernamefromRequestSession error: " + err.Error()))
		return
	}

	var data model.TemplateInfo
	data.CurrentUser = username
	//acc := model.NewAccount()
	//acc.Username = username
	user := model.NewUser()
	user.Username = username
	acclist, err := user.GetAccountList()
	if err != nil {
		log.Println("ReportTemplateHandler acc.GetInfoList() error:", err)
		w.Write([]byte("ReportTemplateHandler error: " + err.Error()))
		return
	}
	data.AccountList = acclist
	data.UsingReport = "yes"
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

	if r.Method == "POST" {
		fmt.Fprintf(w, "U send post request")
	}
}
