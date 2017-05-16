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
	//session, err := store.Get(r, "sessionSSA")
	//if err != nil {
	//	log.Println(err)
	//
	//}

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

	//u := model.NewUser()
	//u.Username = name
	//userinfo, err := u.GetInfo()
	//if err != nil {
	//	log.Println(err)
	//
	//}
	log.Println("IndexHandler getUsernamefromRequestSession name: ", name)
	//log.Printf("User Info: %+v\n", userinfo)
	//objectId := userinfo.Id
	//log.Println(objectId.Hex())
	//Idstring = objectId.Hex()

	// if !IsUserLoggedIn(r) {
	// 	//w.Write([]byte(" you are not logged "))
	// 	http.Redirect(w, r, "/", 302)
	// 	return
	// }
	logged, err := utils.IsUserLoggedIn(r)
	if err != nil {
		if logged {
			data.CurrentUser = name
		} else {

			data.CurrentUser = ""
		}

	} else {
		log.Println("IndexHandler IsUserLoggedIn error: ", err)
		w.Write([]byte("IndexHandler IsUserLoggedIn error" + err.Error()))
		return
	}

	//if session.Values["loggedin"] != nil && session.Values["loggedin"].(string) == "true" {
	//	data.CurrentUser = name
	//} else {
	//	data.CurrentUser = ""
	//}

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
