package vkontakte

import (
	"log"
	"net/http"
)

func GetVKAuthCode(w http.ResponseWriter, r *http.Request) {

	// https://oauth.vk.com/authorize?client_id=1&display=page&redirect_uri=http://example.com/callback&scope=friends&response_type=code
	log.Println(" --- :GetVKAuthCode used ")
	VKurl := "https://oauth.vk.com/authorize?client_id=" + Config.VKAppID +
		"&scope=stats,ads,email&redirect_uri=" + Config.VKRedirectURL + "&response_type=code"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// sending back to client endpoint of redirecting url
	// with id of this app and client's Yandex account login
	w.Write([]byte(VKurl))
	return

}
