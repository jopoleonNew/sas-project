package user

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
)

func ReportTemplateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := r.Context().Value("username").(string)
	accountlogin := vars["accountlogin"]
	logrus.Infof("ReportTemplateHandler used with username: %s, accountlogin: %s", username, accountlogin)
	var data model.TemplateInfo
	a := model.NewAccount2("", "", "", "")
	a.Accountlogin = accountlogin
	info, err := a.GetInfo()
	if err != nil {
		logrus.Errorf("ReportTemplateHandler a.GetInfo() %v error: %v", a, err)
		http.Error(w, fmt.Sprintf("can't find in db Yandex account %v \n error: %+v:", a, err), http.StatusBadRequest)
		return
	}
	data.AccountList = append(data.AccountList, info)
	data.CurrentUser = username
	t, err := template.New("reports.tmpl").ParseFiles(
		"static/templates/header.tmpl",
		"static/templates/reports.tmpl")
	if err != nil {
		logrus.Error("ReportTemplateHandler template.New(reports.tmpl).ParseFiles error: ", err)
		fmt.Fprintf(w, err.Error())
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		logrus.Error("ReportTemplateHandler t.Execute(w, data) error: ", err)
		fmt.Fprintf(w, err.Error())
		return
	}

	if r.Method == "POST" {
		fmt.Fprintf(w, "You've send post request")
		return
	}
}
