package vkontakteAPI

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/sirupsen/logrus"
)

/*
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

func GetAccessToken(appID, appSecret, redirectURL, authURL, code string) (VKtoken, error) {
	//https://oauth.vk.com/access_token?
	//client_id=APP_ID&
	//client_secret=APP_SECRET&
	//code=7a6fa4dff77a228eeda56603b8f53806c883f011c40b72630bb50df056f6479e52a&
	//redirect_uri=REDIRECT_URI&
	var token VKtoken
	VKurl := authURL +
		"?client_id=" + appID +
		"&client_secret=" + appSecret +
		"&code=" + code +
		"&redirect_uri=" + redirectURL
	client := &http.Client{}
	r, err := http.NewRequest("POST", VKurl, nil)
	if err != nil {
		logrus.Errorln("GetAccessToken http.NewRequest error: ", err)
		return token, fmt.Errorf("VkAuthorize http.NewRequest error: %s", err)
	}
	resp, err := client.Do(r)
	if err != nil {
		logrus.Errorln("GetAccessToken client.Do(r) error: ", err)
		return token, fmt.Errorf("VkAuthorize client.Do(r) error: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln("GetAccessToken ioutil.ReadAll(resp.Body) error", err)
		return token, fmt.Errorf("VkAuthorize ioutil.ReadAll(resp.Body) error: %s", err)
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		logrus.Warnf("GetAccessToken bad request repsonse body, trying to unmarshal err %s", err)
		var vkerr VKtokenErr
		err = json.Unmarshal(body, &vkerr)
		if err != nil {
			logrus.Fatalf("response GetAccessToken YandexOauthError json.Unmarshal: \n Indefined body: %s %s", err, string(body))
			return token, err
		}
		return token, errors.New("GetAccessToken VK API error: " + vkerr.Error + " " + vkerr.ErrorDesсription)
	}
	return token, nil
}
