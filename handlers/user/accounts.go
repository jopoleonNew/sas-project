package user

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"strings"

	"gogs.itcloud.pro/SAS-project/sas/app"
	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"
)

func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	//
	//session, err := store.Get(r, "sessionSSA")
	//if err != nil {
	//	log.Println(err)
	//}
	//username, err := utils.GetUsernamefromRequestSession(r)
	//if err != nil {
	//	log.Println("AccountsHandler GetUsernamefromRequestSession err: ", err)
	//	//w.Write([]byte("AccountsHandler GetUsernamefromRequestSession err: " + err.Error()))
	//	http.Redirect(w, r, "/", http.StatusSeeOther)
	//	return
	//}
	//log.Println("AccountsHandler used with username: ", username)
	//if session.Values["loggedin"].(string) == "false" {
	//	http.Redirect(w, r, "/", 302)
	//}
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
	//if session.Values["username"].(string) != username {
	//	w.Write([]byte("You are not logged as " + username))
	//	//http.Redirect(w, r, "/", 302)
	//	return
	//}

	type datas struct {
		UsingReport string
		CurrentUser string
		AccountList []model.Account
	}
	var data datas
	data.CurrentUser = username
	//acc := model.NewAccount()
	//acc.Username = username
	//acclist, err := acc.GetInfoList()
	acclist, err := user.GetAccountList()
	if err != nil {
		log.Println("AccountsHandler user.GetAccountList() error:", err)
		w.Write([]byte("AccountsHandler user.GetAccountList() error: " + err.Error()))
		return
	}
	for _, uaccount := range acclist {
		if uaccount.YandexRole == "agency" {
			acc := model.NewAccount()
			acc.Username = username
			acc.Accountlogin = uaccount.Accountlogin
			agencyInfo, err := acc.GetInfo()
			if err != nil {
				log.Println("AccountsHandler agencyInfo.GetInfo error:", err)
				w.Write([]byte("AccountsHandler agencyInfo.GetInfo error: " + err.Error()))
				return
			}
			for _, agencyAccountLogin := range agencyInfo.AgencyClients {
				agencyAcc := model.NewAccount()
				agencyAcc.Username = username
				agencyAcc.Accountlogin = agencyAccountLogin
				agencyAccInfo, err := agencyAcc.GetInfo()
				if err != nil {
					log.Println("AccountsHandler agencyAccInfo.GetInfo error:", err)
					w.Write([]byte("AccountsHandler agencyAccInfo.GetInfo error: " + err.Error()))
					return
				}
				log.Println("Inside AccountHandler. Agency's AccountInfo : %+v", agencyAccInfo)
				data.AccountList = append(data.AccountList, agencyAccInfo)
			}
		}
	}
	log.Println("AccountsHandler user.GetAccountList(): ", acclist)
	data.AccountList = append(data.AccountList, acclist...)
	//data.AccountList = acclist

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
	acc.Accountlogin = strings.ToLower(r.FormValue("accountlogin"))
	acc.Source = r.FormValue("sourcename")
	acc.YandexRole = r.FormValue("accrole")
	log.Println("AddAccountHandler: ", acc.Accountlogin, acc.Source, acc.YandexRole)

	user := model.NewUser()
	user.Username = username
	userinfo, err := user.GetInfo()
	if err != nil {
		log.Println("AddAccountHandler user.GetInfo() error: ", err)
		return
	}

	for _, a := range user.AccountList {
		if acc.Accountlogin == a {
			log.Println("Аккаунт с таким именем уже существует у пользователя " + username)
			w.Write([]byte("Аккаунт с именем " + acc.Accountlogin + " уже существует у пользователя " + username))
			return
		}
	}
	acc.Email = userinfo.Email
	acc.Status = "notactive"
	acc.SsaAppYandexID = Config.YandexDirectAppID
	acc.SsaAppYandexSecret = Config.YandexDirectAppSecret
	user.AccountList = append(user.AccountList, acc.Accountlogin)
	err = user.AdvanceUpdate()
	if err != nil {
		log.Println("AddAccountHandler user.AdvanceUpdate() error ", err)
		w.Write([]byte("AddAccountHandler user.AdvanceUpdate() error " + err.Error()))
		return
	}
	err = acc.AdvanceUpdate()
	if err != nil {
		log.Println("AddAccountHandler acc.AdvanceUpdate error ", err)
		w.Write([]byte("AddAccountHandler acc.AdvanceUpdate error " + err.Error()))
		return
	}

	//}
	w.Write([]byte("Succsess. Аккаунт добавлен."))
}

func AddVkAccount(w http.ResponseWriter, r *http.Request) {

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
