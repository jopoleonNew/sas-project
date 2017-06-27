package vkontakteAPI

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	API_METHOD_URL    = "https://api.vk.com/method/"
	paramCode         = "code"
	paramToken        = "access_token"
	paramVersion      = "v"
	paramAppID        = "client_id"
	paramScope        = "scope"
	paramRedirectURI  = "redirect_uri"
	paramDisplay      = "display"
	paramHTTPS        = "https"
	paramResponseType = "response_type"

	oauthHost         = "oauth.vk.com"
	oauthDisplay      = "page"
	oauthPath         = "/authorize/"
	oauthResponseType = "token"
	oauthRedirectURI  = "https://oauth.vk.com/blank.html"
	oauthScheme       = "https"

	defaultHost    = "api.vk.com"
	defaultPath    = "/method/"
	defaultScheme  = "https"
	defaultVersion = "5.35"
	defaultMethod  = "GET"
	defaultHTTPS   = "1"

	maxRequestsPerSecond = 3
	minimumRate          = time.Second / maxRequestsPerSecond
	methodExecute        = "execute"
	maxRequestRepeat     = 10
)

type RequestType struct {
	Method string     `json:"method"`
	Token  string     `json:"token"`
	Values url.Values `json:"values"`
}

func Request(token, methodName string, params map[string]string) (string, error) {
	u, err := url.Parse(API_METHOD_URL + methodName)
	if err != nil {
		log.Println("VkAPI http.Get(u.String()) error: ", err)
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	q.Set("access_token", token)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("VkAPI http.Get(u.String()) error: ", err)
		return "", err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("VkAPI http.Get(u.String()) error: ", err)
		return "", err
	}

	return string(content), nil
}

//https://api.vk.com/method/METHOD_NAME?PARAMETERS&access_token=ACCESS_TOKEN&v=V
