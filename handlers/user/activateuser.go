package user

import (
	"log"
	"net/http"

	"gogs.itcloud.pro/SAS-project/sas/model"
)

func ActivateUserHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	log.Println("ActivateUserHandler income URL query: ", r.URL.Query())

	username := query["username"]
	if username == nil || len(username) == 0 {
		log.Println("ActivateUserHandler request received without username")
		w.Write([]byte("ActivateUserHandler request received without username"))
		return
	}
	// "state" is the Yandex account login sent with GetAuthCodeYandexHandler()
	akey := query["activationkey"]
	if akey == nil || len(akey) == 0 {
		log.Println("ActivateUserHandler request received without activation key")
		w.Write([]byte("ActivateUserHandler request received without activation key"))
		return
	}
	user := model.NewUser()
	user.Username = username[0]
	uInfo, err := user.GetInfo()
	if err != nil {
		log.Println("ActivateUserHandler user.GetInfo() error: ", err)
		w.Write([]byte("ActivateUserHandler user.GetInfo() error: " + err.Error()))
		return
	}
	if uInfo.ActivationKey == akey[0] {
		user.IsActivated = "true"
		user.AdvanceUpdate()
		//w.Header().Set("message", "User "+user.Username+" successfully activated")
		http.Redirect(w, r, "/", 301)
	} else {
		w.Write([]byte("Wrong activation key"))
	}
	return
}
