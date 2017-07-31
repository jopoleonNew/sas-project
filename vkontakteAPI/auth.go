package vkontakteAPI

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"fmt"

	"github.com/sirupsen/logrus"
)

/*
VK app ID: 6082545

Защищённый ключ: 8JD0UiJ1KJL43DNm816q

Сервисный ключ доступа:	8c44b2ba8c44b2ba8c44b2ba978c187d4b88c448c44b2bad50dbf1cab30822d4b1745d5

https://oauth.vk.com/authorize?
 client_id=APP_ID&
 scope=SETTINGS&
 redirect_uri=REDIRECT_URI&
 response_type=code
 stats	Доступ к статистике групп и приложений пользователя, администратором которых он является.
ads	Доступ к расширенным методам работы с рекламным API.
offline	Доступ к API в любое время со стороннего сервера.
nohttps	Возможность осуществлять запросы к API без HTTPS.
Внимание, данная возможность находится на этапе тестирования и может быть изменена.
*/
// MakeVKOauthRequest makes request to Yandex Direct Api with confirmation code
// (https://oauth.yandex.ru/token) to get access token
// Обмен кода подтверждения на токен
//https://sas.itcloud.pro/getauthcodevk

type VKtoken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
	Email       string `json:"email"`
}

type VKtokenErr struct {
	Error            string `json:"error"`
	ErrorDesсription string `json:"error_desсription"`
}

var logr = logrus.New()

func VkAccessToken(appID, appSecret, redirectURL, code string) (VKtoken, error) {
	//https://oauth.vk.com/access_token?
	//client_id=APP_ID&
	//client_secret=APP_SECRET&
	//code=7a6fa4dff77a228eeda56603b8f53806c883f011c40b72630bb50df056f6479e52a&
	//redirect_uri=REDIRECT_URI&
	var token VKtoken
	VKurl := "https://oauth.vk.com/access_token?" +
		"client_id=" + appID +
		"&client_secret=" + appSecret +
		"&code=" + code +
		"&redirect_uri=" + redirectURL

	logr.Infoln("VkAccessToken used")

	client := &http.Client{}
	r, err := http.NewRequest("POST", VKurl, nil)
	//r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		logr.Errorln("VkAuthorize http.NewRequest error: ", err)
		return token, fmt.Errorf("VkAuthorize http.NewRequest error: %s", err)
	}
	resp, err := client.Do(r)
	if err != nil {
		logr.Errorln("VkAuthorize client.Do(r) error: ", err)
		return token, fmt.Errorf("VkAuthorize client.Do(r) error: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logr.Errorln("VkAuthorize ioutil.ReadAll(resp.Body) error", err)
		return token, fmt.Errorf("VkAuthorize ioutil.ReadAll(resp.Body) error: %s", err)
	}
	log.Println("VkAuthorize response: ", string(body))
	err = json.Unmarshal(body, &token)
	if err != nil {
		logr.Warnf("VkAuthorize bad request repsonse body, trying to unmarshal err %s", err)
		var vkerr VKtokenErr
		err = json.Unmarshal(body, &vkerr)
		if err != nil {
			logr.Fatalf("response VkAuthorize YandexOauthError json.Unmarshal: \n Indefined body: %s %s", err, string(body))
			return token, err
		}
		return token, errors.New("VkAccessToken VK API error: " + vkerr.Error + " " + vkerr.ErrorDesсription)
	}
	//log.Println("////\n\n TOKEN FROM VKONTAKTE: ", token)
	return token, nil
}
