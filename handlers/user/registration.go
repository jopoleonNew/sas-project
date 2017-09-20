package user

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("static/templates/signup.html")
	if err != nil {
		logrus.Error("SignUpHandler template parse error: ", err)
		fmt.Fprintf(w, err.Error())
		return
	}
	t.ExecuteTemplate(w, "signup.html", nil)
}

func SignUpSubmitHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		logrus.Error("SignUpSubmitHandler error: ", err)
		w.Write([]byte("Registration server error " + err.Error()))
		return
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	name := r.FormValue("name")
	organization := r.FormValue("organization")

	u := model.NewUser()
	u.Username = username
	//u.Username = username
	u.Password = password
	u.Email = email
	u.Name = name
	u.Organization = organization
	u.IsActivated = "false"

	if username != "" && password != "" && email != "" {
		if model.ValidateEmail(email) {
			yes, err := u.IsExist()
			if err != nil {
				logrus.Error("SignUpSubmitHandler u.IsExists error: ", err)
				w.Write([]byte("Registration server error " + err.Error()))
				return
			}
			if !yes {
				akey := model.RandStringBytes(32)
				u.ActivationKey = akey

				err := u.Create()
				if err != nil {
					logrus.Error("SignUpSubmitHandler u.Update() error: ", err)
					w.Write([]byte("Registration server error " + err.Error()))
					return
				}
				session.Values["loggedin"] = "false"
				session.Values["username"] = username
				session.Save(r, w)
				err = model.SendEmailwithKey(username, akey, email, r.Host)
				if err != nil {
					logrus.Error("SignUpSubmitHandler SendEmailwithKey error: ", err)
					w.Write([]byte("Registration server error " + err.Error()))
					return
				}
				w.Write([]byte("Registration is successful. Check your email for activation letter."))
				return
			}
			w.Write([]byte("User with such login or email already exists"))
			return
		}
		w.Write([]byte("Email is invalid format"))
		return
	}
	w.Write([]byte("Some of registration fields are empty!"))
}
