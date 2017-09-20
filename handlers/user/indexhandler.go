package user

import (
	"fmt"
	"html/template"
	"log"

	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/utils"

	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	//
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println(err)

	}

	var data model.TemplateInfoStruct

	//log.Println("IndexHandler session.Values", session.Values["loggedin"])
	name, err := utils.GetUsernamefromRequestSession(r)
	if err != nil {
		if err == utils.ErrorEmptySessionUsername {
			//log.Println("Inside check")
			data.CurrentUser = ""
		} else {
			log.Println("IndexHandler GetUsernamefromRequestSession error: ", err)
			w.Write([]byte("IndexHandler GetUsernamefromRequestSession error" + err.Error()))
			return
		}
	}
	loggedin, err := IsUserLoggedIn(r)
	if err != nil {
		log.Println("IndexHandler IsUserLoggedIn error: ", err)
		w.Write([]byte("IndexHandler IsUserLoggedIn error" + err.Error()))
		return
	} else {
		if loggedin {
			data.CurrentUser = name
		} else {
			data.CurrentUser = ""
		}
	}

	if session.Values["loggedin"] != nil && session.Values["loggedin"].(string) == "true" {
		data.CurrentUser = name
	} else {
		data.CurrentUser = ""
	}

	log.Println("IndexHandler getUsernamefromRequestSession name: ", name)

	t := template.Must(template.New("index.tmpl").ParseFiles(
		"static/templates/index.tmpl",
		"static/templates/header.tmpl",
		"static/templates/login.tmpl"))

	err = t.ExecuteTemplate(w, "index.tmpl", &data)
	if err != nil {
		log.Println("IndexHandler error: ", err)
		fmt.Fprintf(w, err.Error())
	}

}

func IsUserLoggedIn(r *http.Request) (bool, error) {
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println(err)

	}
	log.Println("IsUserLoggedIn values: ", session.Values)

	if session.Values["loggedin"] != nil &&
		session.Values["loggedin"].(string) == "true" &&
		len(session.Values) == 0 {
		return true, nil
	} else {
		//log.Println("IsUserLoggedIn: user not loggedin: \n RemoteAddres:", r.RemoteAddr, " \n RequestURI: ", r.RequestURI)
		return false, nil
	}

}
