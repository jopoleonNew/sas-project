package yad

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// MakeYandexOauthRequest makes request to Yandex Direct Api with confirmation code
// (https://oauth.yandex.ru/token) to get access token
// Обмен кода подтверждения на токен
func MakeYandexOauthRequest(code string) (YandexTokenbody, error) {
	var token YandexTokenbody
	yandextokenurl := "https://oauth.yandex.ru/token"
	log.Println("MakeYandexOauthRequest used")
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("client_id", application.ID)
	form.Add("client_secret", application.Secret)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", yandextokenurl, strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		log.Println("MakeYandexOauthRequest client.Do(r) error: ", err)
		return token, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("MakeYandexOauthRequest ioutil.ReadAll(resp.Body) error", err)
		return token, err
	}
	log.Println("MakeYandexOauthRequest response: ", string(body))
	if string(body) == `{"error_description": "Invalid code", "error": "bad_verification_code"}` {
		log.Println("MakeYandexOauthRequest response: .............................")

		errcode := errors.New(`{"error_description": "Invalid code", "error": "bad_verification_code"}`)
		return token, errcode
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Println("MakeYandexOauthRequest bad request repsonse body, trying to unmarshal err ", err)
		var yaderr YandexV5Error
		err = json.Unmarshal(body, &yaderr)
		if err != nil {
			log.Println("response MakeYandexOauthRequest YandexOauthError json.Unmarshal: \n Indefined body: ", err, string(body))
			return token, err
		}
		return token, errors.New("YandexDirectAPI error: " + yaderr.ErrorString + " " + yaderr.ErrorDescription)
	}
	return token, nil
}
