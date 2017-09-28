package user

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"encoding/json"

	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"
)

// AccountsHandler returns to client the list of client's account from all sources
func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	//
	username := r.Context().Value("username").(string)
	log.Println("AccountsHandler used with username: ", username)
	user := model.NewUser()
	user.Username = username
	exist, err := user.IsExist()
	if err != nil {
		log.Println("AccountsHandler u.IsExist() error:", err)
		w.Write([]byte("AccountsHandler u.IsExist() error " + err.Error()))
		return
	}

	if !exist {
		w.Write([]byte("No such user found"))
		//http.Redirect(w, r, "/", 302)
		return
	}

	type datas struct {
		UsingReport string
		CurrentUser string
		AccountList []model.Account2
	}
	var data datas
	data.CurrentUser = username
	acc := model.NewAccount2("", "", "", "")
	acc.Owners = []string{username}
	acclist, err := acc.GetAccountList()
	if err != nil {
		log.Println("AccountsHandler acc.GetAccountList() error:", err)
		w.Write([]byte("AccountsHandler acc.GetAccountList() error: " + err.Error()))
		return
	}
	//logrus.Infof("Inside AccountHandler acc.GetAccountList() res: %+v", acclist)
	data.AccountList = append(data.AccountList, acclist...)
	t, err := template.New("accounts_new.tmpl").ParseFiles(
		"static/templates/header.tmpl",
		"static/templates/accounts_new.tmpl",
		"static/templates/login.tmpl")
	if err != nil {
		log.Println("AccountsHandler template.New error: ", err)
		fmt.Fprintf(w, err.Error())
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Println("AccountsHandler t.Execute error: ", err)
		fmt.Fprintf(w, err.Error())
	}

}
func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteAccountHandler used")
	username, err := utils.GetUsernamefromRequestSession(r)
	if err != nil {
		log.Println("DeleteAccountHandler GetUsernamefromRequestSession error: ", err)
		w.Write([]byte("DeleteAccountHandler GetUsernamefromRequestSession error: " + err.Error()))
		return
	}
	r.ParseForm()
	accountlogin := r.FormValue("accountlogin")
	user := model.NewUser()
	//acc := model.NewAccount()
	user.Username = username
	user.AccountList = append(user.AccountList, accountlogin)
	err = user.RemoveAccount()
	if err != nil {
		log.Println("DeleteAccountHandler acc.Remove() error: ", err)
		w.Write([]byte("DeleteAccountHandler acc.Remove() error: " + err.Error()))
		return
	}
}
func AccountsData(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	log.Println("AccountsData used with username: ", username)
	user := model.NewUser()
	user.Username = username
	exist, err := user.IsExist()
	if err != nil {
		log.Println("AccountsData u.IsExist() error:", err)
		w.Write([]byte("AccountsData u.IsExist() error " + err.Error()))
		return
	}

	if !exist {
		w.Write([]byte("No such user found"))
		//http.Redirect(w, r, "/", 302)
		return
	}

	var AccountList []model.Account2
	acc := model.NewAccount2("", "", "", "")
	acc.Owners = []string{username}
	AccountList, err = acc.GetAccountList()
	if err != nil {
		log.Println("AccountsData acc.GetAccountList() error:", err)
		w.Write([]byte("AccountsData acc.GetAccountList() error: " + err.Error()))
		return
	}
	b, err := json.Marshal(AccountList)
	if err != nil {
		log.Println("AccountsData json.Marshal(AccountList) error:", err)
		w.Write([]byte("AccountsData json.Marshal(AccountList) error: " + err.Error()))
		return
	}
	w.Write(b)

}
