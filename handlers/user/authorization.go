package user

import (
	"log"

	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"

	"net/http"
)


func LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println("LoginSubmitHandler store.Get error:", err)
		w.Write([]byte("LoginSubmitHandler store.Get error " + err.Error()))
		return
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Println("LoginSubmitHandler: ", username, password)
	user := model.NewUser()
	user.Username = username
	user.Password = password
	exist, err := user.IsExist()
	if err != nil {
		log.Println("LoginSubmitHandler u.IsExist()error:", err)
		w.Write([]byte("LoginSubmitHandler u.IsExist() error " + err.Error()))
		return
	}
	if !exist {
		log.Println("no such user found")
		w.Write([]byte("No such user found"))
		return
	}

	valid, err := user.IsPasswordValid(password)
	if err != nil {
		log.Println("LoginSubmitHandler u.IsPasswordValid error: ", err)
		w.Write([]byte("User validation error: " + err.Error()))
		return
	}
	if !valid {
		log.Println("Password not valid")
		w.Write([]byte("Incorrect password"))
		return
	}

	session.Values["username"] = username
	//session.Values["password"] = password
	session.Values["loggedin"] = "true"
	session.Save(r, w)
	w.Write([]byte("Success!"))
	//http.Redirect(w, r, "/", 302)
	return

}
func LogoutSubmitHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println("LogoutSubmitHandler store.Get error: ", err)
		w.Write([]byte("LogoutSubmitHandler store.Get error" + err.Error()))
		return
	}
	username, err := utils.GetUsernamefromRequestSession(r)
	if err != nil {
		log.Println("LogoutSubmitHandler getUsernamefromRequestSession error", err)
		w.Write([]byte("LogoutSubmitHandler getUsernamefromRequestSession error" + err.Error()))
		return
	}
	log.Println("LogoutSubmitHandler requset username: ", username)
	//log.Println("Inside Logout. Is user registred: ", BDctl.IsUserRegistered(username),
	//	"session.Values[loggedin]: ", session.Values["loggedin"])
	user := model.NewUser()
	user.Username = username
	//u.Password = password
	exist, err := user.IsExist()
	if err != nil {
		log.Println("LoginSubmitHandler u.IsExist()error:", err)
		w.Write([]byte("LoginSubmitHandler u.IsExist() error " + err.Error()))
		return
	}
	if exist && session.Values["loggedin"].(string) == "true" {
		session.Options.MaxAge = -1
		session.Values["username"] = ""
		session.Values["password"] = ""
		session.Values["loggedin"] = "false"

	}
	//var options sessions.Options
	//options.MaxAge = -1
	//session.
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
	return

}
