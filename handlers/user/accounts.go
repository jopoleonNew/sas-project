package user

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"gogs.itcloud.pro/SAS-project/sas/app"
	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"
)

func AccountsHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println(err)
	}
	username, err := utils.GetUsernamefromRequestSession(r)
	if err != nil {
		log.Println("AccountsHandler GetUsernamefromRequestSession err: ", err)
		//w.Write([]byte("AccountsHandler GetUsernamefromRequestSession err: " + err.Error()))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	log.Println("AccountsHandler used with username: ", username)
	if session.Values["loggedin"].(string) == "false" {
		http.Redirect(w, r, "/", 302)
	}
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
	if session.Values["username"].(string) != username {
		w.Write([]byte("You are not logged as " + username))
		//http.Redirect(w, r, "/", 302)
		return
	}

	type datas struct {
		UsingReport string
		CurrentUser string
		AccountList []model.Account
	}
	var data datas
	data.CurrentUser = username
	acc := model.NewAccount()
	acc.Username = username
	acclist, err := acc.GetInfoList()
	if err != nil {
		log.Println("AccountsHandler acc.GetInfoList() error:", err)
		w.Write([]byte("AccountsHandler error: " + err.Error()))
		return
	}
	data.AccountList = acclist

	t, err := template.New("accounts2.tmpl").ParseFiles(
		"static/templates/header.tmpl",
		"static/templates/accounts2.tmpl",
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

func AddAccountHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AddAccountHandler used")
	Config = app.GetConfig()
	username, err := utils.GetUsernamefromRequestSession(r)
	if err != nil {
		log.Println("AddAccountHandler GetUsernamefromRequestSession err: ", err)
		w.Write([]byte("AddAccountHandler GetUsernamefromRequestSession err: " + err.Error()))
		return
	}

	//accountlogin := r.FormValue("accountlogin")
	//sourcename := r.FormValue("sourcename")
	//accrole := r.FormValue("accrole")

	err = r.ParseForm()
	if err != nil {
		log.Println("AddAccountHandler r.ParseForm() err: ", err)
		w.Write([]byte("AddAccountHandler r.ParseForm() err: " + err.Error()))
		return
	}

	acc := model.NewAccount()
	acc.Username = username
	acc.Accountlogin = r.FormValue("accountlogin")
	acc.Source = r.FormValue("sourcename")
	acc.YandexRole = r.FormValue("accrole")
	log.Println("AddAccountHandler: ", acc.Accountlogin, acc.Source, acc.YandexRole)

	user := model.NewUser()
	user.Username = username
	userinfo, err := user.GetInfo()
	if err != nil {
		log.Println("AccountsHandler user.GetInfo() error: ", err)
		return
	}
	//func (ctl *Controller) IsAccountUnique(username, accountlogin string) (bool, error)

	exists, err := acc.IsExist()
	if err != nil && err != model.ErrAccNotFound {
		log.Println("AccountsHandler acc.IsExist() error: ", err)
		//w.Write([]byte("Аккаунт у этого пользователя с таким именем уже существует"))
		return
	}

	if exists {
		w.Write([]byte("Аккаунт с именем " + acc.Accountlogin + " уже есть в базе. "))
		return
	} else {
		acc.Email = userinfo.Email
		acc.Status = "notactive"
		acc.SsaAppYandexID = Config.YandexDirectAppID
		acc.SsaAppYandexSecret = Config.YandexDirectAppSecret
		err := acc.AdvanceUpdate()
		if err != nil {
			log.Println("AddAccountToDB error ", err)
			w.Write([]byte("AddAccountToDB error " + err.Error()))
			return

		}
	}
	w.Write([]byte("Succsess. Аккаунт добавлен."))
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {

	username, err := utils.GetUsernamefromRequestSession(r)
	if err != nil {
		log.Println("DeleteAccountHandler GetUsernamefromRequestSession error: ", err)
		w.Write([]byte("DeleteAccountHandler GetUsernamefromRequestSession error: " + err.Error()))
		return
	}
	r.ParseForm()
	accountlogin := r.FormValue("accountlogin")
	acc := model.NewAccount()
	acc.Username = username
	acc.Accountlogin = accountlogin
	err = acc.Remove()
	if err != nil {
		log.Println("DeleteAccountHandler acc.Remove() error: ", err)
		w.Write([]byte("DeleteAccountHandler acc.Remove() error: " + err.Error()))
		return
	}
}
