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

	u := model.NewUser()
	u.Username = username
	//u.Username = username
	u.Password = password
	u.Email = email
	u.Name = name
	u.Organization = organization

	if username != "" && password != "" && email != "" {
		if model.ValidateEmail(email) {
			yes, err := u.IsExist()
			if err != nil {
				log.Println("SignUpSubmitHandler u.IsExists error: ", err)
				w.Write([]byte("Registration server error " + err.Error()))
				return
			}
			if !yes {
				err := u.Update()
				if err != nil {
					log.Println("SignUpSubmitHandler u.Update() error: ", err)
					w.Write([]byte("Registration server error " + err.Error()))
					return
				}
				session.Values["registered"] = "true"
				session.Values["loggedin"] = "false"
				session.Values["username"] = username
				session.Values["password"] = password
				session.Values["email"] = email
				session.Values["name"] = name
				session.Values["organization"] = organization
				session.Save(r, w)
				w.Write([]byte("Registration succsessfull "))
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
