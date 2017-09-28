package user

import (
	"log"
	"net/http"

	"strings"

	"io/ioutil"

	"encoding/json"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/shared/config"
)

func LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("LoginSubmitHandler used:\n %+v ", r)
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Origin, Content-Type, X-Auth-Token, Authorization, Username, Password")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		//w.Header().Set("Access-Control-Allow-Origin", "https://"+r.Host+"/loginsubmit"+", http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.WriteHeader(200)
		return
	}
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		logrus.Error("LoginSubmitHandler store.Get error: ", err)
		http.Error(w, "Internal server error: "+err.Error(), 500)
		return
	}
	r.ParseForm()
	//logrus.Info("LoginSubmitHandler username: ", username, ", password: ", password)
	username := strings.ToLower(r.FormValue("username"))
	password := r.FormValue("password")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Origin, Content-Type, X-Auth-Token, Authorization, Username, Password")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Access-Control-Allow-Origin", "https://"+r.Host+"/loginsubmit"+", http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	logrus.Infoln("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	if username == "" || password == "" {
		{
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logrus.Error("LoginSubmitHandler ioutil.ReadAll(r.Body) error:", err)
				http.Error(w, "Internal server error: error: "+err.Error(), 500)
				return
			}
			type authInfo struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			var auth authInfo
			err = json.Unmarshal(body, &auth)
			if err != nil {
				logrus.Error("LoginSubmitHandler json.Unmarshal(body,&auth) error:", err)
				http.Error(w, "Internal server error: "+err.Error(), 500)
				return
			}
			defer r.Body.Close()
			if auth.Username == "" || auth.Password == "" {
				logrus.Error("LoginSubmitHandler username or password are empty in request")
				http.Error(w, "LoginSubmitHandler username or password are empty in request.", 400)
				return
			}
			username = strings.ToLower(auth.Username)
			password = auth.Password
		}

	}
	//logrus.Info("LoginSubmitHandler username: ", username, ", password: ", password)
	username = strings.ToLower(username)
	user := model.NewUser()
	user.Username = username
	user.Password = password
	exist, err := user.IsExist()
	if err != nil {
		log.Println("LoginSubmitHandler u.IsExist()error:", err)
		http.Error(w, "Internal server error: "+err.Error(), 500)
		return
	}
	if !exist {
		//logrus.Println("no such user found")
		http.Error(w, "No such user found", 404)
		return
	}

	valid, err := user.IsPasswordValid(password)
	if err != nil {
		logrus.Error("LoginSubmitHandler u.IsPasswordValid error: ", err)
		http.Error(w, "Internal server error: "+err.Error(), 500)
		return
	}
	if !valid {
		http.Error(w, "Password incorrect", 400)
		return
	}
	uInfo, err := user.GetInfo()
	if uInfo.IsActivated == "false" {
		log.Println("User ", username, " is not activated.")
		http.Error(w, "Пользователь "+username+" не активирован. Проверьте свой почтовый ящик, там должно быть активационное письмо.", 400)
		return
	}
	session.Values["username"] = username
	session.Values["loggedin"] = "true"
	session.Save(r, w)
	w.WriteHeader(200)
	w.Write([]byte("Success!"))
	return
}

func LogoutSubmitHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println("LogoutSubmitHandler store.Get error: ", err)
		w.Write([]byte("LogoutSubmitHandler store.Get error" + err.Error()))
		return
	}
	username := r.Context().Value("username").(string)
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
		//session.Options.MaxAge = -1
		session.Values["username"] = ""
		//session.Values["password"] = ""
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

func IsLoggedIn(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte(config.GetConfig().SessionSecret))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Origin, Content-Type, X-Auth-Token, Authorization, Username, Password")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Access-Control-Allow-Origin", "https://"+r.Host+"/loginsubmit"+", http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		logrus.Error("IsLoggedIn store.Get err: ", err.Error())
		http.Error(w, "Internal server erro: IsLoggedIn store.Get err: "+err.Error(), 500)
		return
	}
	//log.Println("CheckIsUserLogged middleware values: ", session.Values)

	if session.Values["loggedin"] != nil && session.Values["loggedin"].(string) == "true" &&
		len(session.Values) != 0 {
		//Add data to context
		//ctx := context.WithValue(r.Context(), "username", session.Values["username"])
		username, ok := session.Values["username"].(string)
		if !ok {
			logrus.Error("IsLoggedIn store.Get err: session.Values[username] type assertion error: ", session.Values["username"])
			http.Error(w, "Internal server error: session.Values[username] type assertion error", 500)
			return
		}
		w.Write([]byte(username))
		//w.WriteHeader(200)
		//next.ServeHTTP(w, r)
	} else {
		http.Error(w, "You are not logged in.", 401)
		return
	}
}
