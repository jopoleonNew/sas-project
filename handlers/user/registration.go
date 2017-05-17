package user

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"gogs.itcloud.pro/SAS-project/sas/model"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("static/templates/signup.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "signup.html", nil)
}

func SignUpSubmitHandler(w http.ResponseWriter, r *http.Request) {
	//if r.Method == http.MethodGet {
	//	t, err := template.ParseFiles("static/templates/signup.html")
	//	if err != nil {
	//		fmt.Fprintf(w, err.Error())
	//	}
	//	t.ExecuteTemplate(w, "signup.html", nil)
	//}

	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println("SignUpSubmitHandler error: ", err)
		w.Write([]byte("Registration server error " + err.Error()))
		return
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	name := r.FormValue("name")
	organization := r.FormValue("organization")
	//log.Println("SignUpHandler check r.RemoteAddr: ", r.RemoteAddr)
	//log.Println("SignUpHandler check r.Host: ", r.Host)
	//log.Println("SignUpHandler check r.RequestURI: ", r.RequestURI)
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
				log.Println("SignUpSubmitHandler u.IsExists error: ", err)
				w.Write([]byte("Registration server error " + err.Error()))
				return
			}
			if !yes {
				akey := model.RandStringBytes(32)
				u.ActivationKey = akey

				err := u.Update()
				if err != nil {
					log.Println("SignUpSubmitHandler u.Update() error: ", err)
					w.Write([]byte("Registration server error " + err.Error()))
					return
				}
				//session.Values["registered"] = "true"
				session.Values["loggedin"] = "false"
				session.Values["username"] = username
				//session.Values["password"] = password
				//session.Values["email"] = email
				//session.Values["name"] = name
				//session.Values["organization"] = organization
				session.Save(r, w)
				err = model.SendEmailwithKey(username, akey, email, r.Host)
				if err != nil {
					log.Println("SignUpSubmitHandler SendEmailwithKey error: ", err)
					w.Write([]byte("Registration server error " + err.Error()))
					return
				}
				//
				//TODO: Sending email with activation code to user's email

				//
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
