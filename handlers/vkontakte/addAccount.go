package vkontakte

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	vk "gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
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
func VKauthorize(w http.ResponseWriter, r *http.Request) {
	//REDIRECT_URI?code=7a6fa4dff77a228eeda56603b8f53806c883f011c40b72630bb50df056f6479e52a
	r.ParseForm()
	query := r.URL.Query()
	log.Println("GetYandexAccessToken income URL query: ", r.URL.Query())

	code := query["code"]
	if code != nil || len(code) != 0 {
		vktoken, err := vk.VkAccessToken(Config.VKAppID, Config.VKAppSecret, Config.VKRedirectURL, code[0])
		if err != nil {
			log.Println("VKauthorize vk.VkAccessToken error: ", err)
			return
		}

		//42f17cfb678d3008ad04df046815c5fdfa3663d984771b92db47955675f7a224c1f259b125062ecfdb04b
		tempToken := "42f17cfb678d3008ad04df046815c5fdfa3663d984771b92db47955675f7a224c1f259b125062ecfdb04b"

		log.Println("Inside VKauthorize vk.VkAccessToken result:::::: ", vktoken)
		response, err := vk.Request(tempToken, "ads.getAccounts", nil)
		if err != nil {
			log.Println("VKauthorize vk.Request error: ", err)
			return
		}
		log.Println("vk.Request ads.getAccounts result: ", response)
		return

	}
	log.Println("Request from Vkontakte received without code")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("VkAuthorize ioutil.ReadAll(resp.Body) error", err)
		w.Write([]byte("VkAuthorize ioutil.ReadAll(resp.Body) error" + err.Error()))
		return
	}
	var token vk.VKtoken

	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Println("VkAuthorize bad request repsonse body, trying to unmarshal err ", err)
		var vkerr vk.VKtokenErr
		err = json.Unmarshal(body, &vkerr)
		if err != nil {
			log.Println("response VkAuthorize VKtokenErr json.Unmarshal: \n Indefined body: ", err, string(body))
			w.Write([]byte("response VkAuthorize VKtokenErr json.Unmarshal: \n Indefined body:" + err.Error()))
			return
		}
		w.Write([]byte("YandexDirectAPI error: " + vkerr.Error + " " + vkerr.ErrorDesсription))
		return
	}
	log.Println("Inside VKauthorize vk.VkAccessToken result:::::: ", token)
}
