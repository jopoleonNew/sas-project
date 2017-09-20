package user

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"gogs.itcloud.pro/SAS-project/sas/model"
	"golang.org/x/crypto/bcrypt"
)

func ForgetPassHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	log.Println("ForgetPassHandler used")
	if r.Method == "GET" {
		t, err := template.ParseFiles("static/templates/forgetpass.tmpl")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		t.ExecuteTemplate(w, "forgetpass.tmpl", nil)
	}

	if r.Method == "POST" {
		r.ParseForm()
		email := r.FormValue("email")
		log.Println("ForgetPassHandler incoming email ", email)
		user := model.NewUser()
		//user.Username = "blank"
		user.Email = email
		exsit, err := user.IsExist()
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		if !exsit {
			log.Println("ForgetPassHandler user not found")
			w.Write([]byte("User with such email not found"))
			return
		}

		key := model.RandStringBytes(64)
		err = model.AddLinkKey(email, key)
		if err != nil {
			log.Println("ForgetPassHandler  model.AddLinkKey(email, key) error: ", err)
			fmt.Fprintf(w, err.Error())
			return
		}
		err = model.SendEmailRestorePass(key, email, r.Host)
		if err != nil {
			log.Println("ForgetPassHandler model.SendEmailRestorePass(key, email, r.Host) error: ", err)
			fmt.Fprintf(w, err.Error())
			return
		}
		w.Write([]byte("Письмо для восстановления паролья отправленно на Ваш почтовый ящик."))
		return
		//password := r.FormValue("password")
	}
}
func RestorePasswordHadler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("static/templates/restorepass.tmpl")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		t.ExecuteTemplate(w, "restorepass.tmpl", nil)
	}
	if r.Method == "POST" {
		query := r.URL.Query()
		log.Println("RestorePasswordHadler income URL query: ", r.URL.Query())
		//changepass?secretkey=

		if query["secretkey"] == nil || len(query["secretkey"]) == 0 {
			log.Println("RestorePasswordHadler request received without username")
			w.Write([]byte("RestorePasswordHadler request received without username"))
			return
		}
		secretkey := query["secretkey"][0]
		email, correct, err := model.MatchKey(secretkey)
		if err != nil {
			log.Println("RestorePasswordHadler model.MatchKey(secretkey) error: ", err)
			w.Write([]byte("RestorePasswordHadler model.MatchKey(secretkey) error: " + err.Error()))
			return
		}
		if correct {
			r.ParseForm()
			pass1 := r.FormValue("pass1")
			pass2 := r.FormValue("pass2")
			log.Println("RestorePasswordHadler recieved passwords: ", pass1, pass2)
			if pass1 == "" || pass2 == "" {
				w.Write([]byte("Пароль не может быть пустым"))
				return
			}
			if pass1 != pass2 {
				w.Write([]byte("Пароль и подтверждение паролья не одинаковы"))
				return
			}
			if len(pass1) < 6 {
				w.Write([]byte("Пароль не может быть короче 6 символов"))
				return
			}
			user := model.NewUser()
			user.Email = email
			passbyte := []byte(pass1)
			hashedPassword, err := bcrypt.GenerateFromPassword(passbyte, bcrypt.DefaultCost)
			if err != nil {
				log.Println("RestorePasswordHadler bcrypt.GenerateFromPassword error: ", err)
				w.Write([]byte("RestorePasswordHadler bcrypt.GenerateFromPassword error: " + err.Error()))
			}
			//u.Username = strings.ToLower(u.Username)
			user.Salt = string(hashedPassword)
			user.Password = pass1
			err = user.Update()
			if err != nil {
				log.Println("RestorePasswordHadler user.Update() error: ", err)
				w.Write([]byte("RestorePasswordHadler user.Update() error: " + err.Error()))
			}
			w.Write([]byte("Пароль успешно изменен"))
			return
		} else {
			w.Write([]byte("Неправильный код в ссылке, либо прошло больше 5 часов после его выдачи вам. Попробуйте еще раз получить письмо с новым кодом"))
		}

	}

}
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {

}
