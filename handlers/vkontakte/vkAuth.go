package vkontakte

import (
	"net/http"
)

// GetVKAuthLink writes to ResponseWriter the VK API Auth Link
// which front-end uses to redirect client to give access to his VK ads
func GetVKAuthLink(w http.ResponseWriter, r *http.Request) {
	//https://oauth.vk.com/authorize?client_id=5490057&display=page&redirect_uri=https://oauth.vk.com/blank.html&scope=friends&response_type=token&v=5.52
	VKurl := "https://oauth.vk.com/authorize?client_id=" + Config.VKAppID +
		"&display=page&scope=offline,stats,ads,email&redirect_uri=" + Config.VKRedirectURL + "&response_type=code"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Write([]byte(VKurl))
}
