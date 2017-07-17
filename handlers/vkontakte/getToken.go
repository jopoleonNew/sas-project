package vkontakte

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	vk "gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
)

func GetVKToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	query := r.URL.Query()
	log.Println("GetYandexAccessToken income URL query: ", r.URL.Query())

	code := query["code"]
	if code == nil || len(code) == 0 {
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