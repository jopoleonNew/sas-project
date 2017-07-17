package yandex

import (
	"log"
	"net/http"
)

// GetAuthCodeYandexHandler handling requests from client (mostly from browser's JS) and returns in response
// the ID of this app and client's Yandex login, so the JS of client could redirect
// him to "https://oauth.yandex.ru/authorize?response_type=code&client_id=&state=" + data
// where :
//      "state" is the client's Yandex login
//      "client_id" is the ID of this app
// /getauthcodeyandex endpoint
// TODO: This functions seems to me like a kludge (костыль), so it really needs to be refactored or even embedded in other functions and deleted

func GetAuthCodeYandexHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	username := r.FormValue("username")
	accountlogin := r.FormValue("accountlogin")
	log.Println("GetAuthCodeYandex used ", username, accountlogin)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// sending back to client endpoint of redirecting url
	// with id of this app and client's Yandex account login
	w.Write([]byte(Config.YandexDirectAppID + "&state=" + accountlogin))
	return

}
