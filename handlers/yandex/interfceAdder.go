// Package yandex provides http handlers for work with Yandex Direct accounts.
package yandex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"log"

	"github.com/nk2ge5k/goyad"
	"github.com/nk2ge5k/goyad/ads"
	"github.com/nk2ge5k/goyad/agencyclients"
	"github.com/nk2ge5k/goyad/campaigns"
	"github.com/nk2ge5k/goyad/clients"
	"github.com/nk2ge5k/goyad/gc"
	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
	yad "gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

type YandexAccountAdder interface {
	ParseURL(url *url.URL) (map[string]string, error)
	GetToken(code string) (yad.YandexTokenbody, error)
	CollectAccountandAddtoBD() (info CreateInfo, err error)
	CollectAgencyClients() (res agencyclients.GetResponse, err error)
	CollectClientInfo() (res clients.GetResponse, err error)
	CollectCampaigns() (res campaigns.GetResponse, err error)
	CollectAds() (res ads.GetResponse, err error)
	AddAccToDB(a model.Account2) error
}

type yclient struct {
	goyad.Client
	AuthURL     string
	RedirectURL string
	ApiURL      string
	AppID       string
	AppSecret   string
	Creator     string
}

func (c *yclient) AddYandexAccount(w http.ResponseWriter, r *http.Request) {
	urlValues, err := c.ParseURL(r.URL)
	if err != nil {
		logrus.Error("Request from YandexOuath received without values")
		http.Error(w, fmt.Sprintf("Request from Yandex received without values %s", urlValues), http.StatusBadRequest)
		return
	}

	code := urlValues["code"]
	token, err := c.GetToken(code)
	if err != nil {
		logrus.Println("AddYandexAccount GetToken error: ", err)
		http.Error(w, fmt.Sprintf("cant get auth token from Yandex with code: %s, error: %v", code, err), http.StatusBadRequest)
		return
	}
	creator := r.Context().Value("username").(string)
	if creator == "" {
		logrus.Errorf("AddYandexAccount r.Context().Value(username) is empty: ", creator)
		http.Error(w, fmt.Sprintf("Can't identify username inside request context: %s", creator), http.StatusBadRequest)
		return
	}
	c.Creator = creator
	c.Login = urlValues["accountlogin"]
	c.Token = goyad.Token{Value: token.AccessToken}
	_, err = c.CollectAccountandAddtoBD()
	if err != nil {
		logrus.Errorf("AddYandexAccount r.Context().Value(username) is empty: ", creator)
		http.Error(w, fmt.Sprintf("Can't identify username inside request context: %s", creator), http.StatusBadRequest)
		return
	}
}
func (c *yclient) ParseURL(url *url.URL) (map[string]string, error) {

	p := make(map[string]string)
	query := url.Query()

	if len(query) == 0 {
		return p, fmt.Errorf("Request received without URL Query: %s", query)
	}
	s := query["code"]
	if s == nil || len(s) == 0 {
		return p, fmt.Errorf("Request received without code %s", query)
	}
	p["code"] = s[0]
	al := query["state"]
	if al == nil || len(al) == 0 {
		return p, fmt.Errorf("Request received without state %s", query)
	}
	p["accountlogin"] = al[0]

	return p, nil
}
func (c *yclient) GetToken(code string) (yad.YandexTokenbody, error) {
	var token yad.YandexTokenbody
	if code == "" {
		return token, fmt.Errorf("code can't be empty")
	}
	//log.Println("GetYandexToken used")
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("client_id", c.AppID)
	form.Add("client_secret", c.AppSecret)

	client := &http.Client{}
	r, err := http.NewRequest("POST", c.AuthURL, strings.NewReader(form.Encode()))
	if err != nil {
		logrus.Error("GetYandexToken http.NewRequest error: ", err)
		return token, err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	logrus.Info("Inside GetToken url ", r.URL)
	resp, err := client.Do(r)
	if err != nil {
		logrus.Error("GetYandexToken client.Do(r) error: ", err)
		return token, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("GetYandexToken ioutil.ReadAll(resp.Body) error", err)
		return token, err
	}
	log.Println("GetYandexToken response: ", string(body))
	if string(body) == `{"error_description": "Invalid code", "error": "bad_verification_code"}` {
		logrus.Error("Error from yandexDirect API: Invalid code, bad_verification_code: ", string(body))
		return token, fmt.Errorf("Error from yandexDirect API: Invalid code, bad_verification_code")
	}
	if strings.Contains(string(body), "error") {
		logrus.Warn("GetYandexToken bad request repsonse body, trying to unmarshal err ", err)
		var yaderr yad.YandexV5Error
		err = json.Unmarshal(body, &yaderr)
		if err != nil {
			logrus.Errorf("response GetYandexToken YandexOauthError json.Unmarshal error: %v \n Indefined body: %s", err, string(body))
			return token, err
		}
		logrus.Error("YandexDirectAPI error: " + yaderr.ErrorString + " " + yaderr.ErrorDescription)
		return token, fmt.Errorf("YandexDirectAPI error: %v, %v", yaderr.ErrorString, yaderr.ErrorDescription)
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		logrus.Errorf("GetToken cant unmarshl body of YandexDirectAPI response: %s, error: %v ", string(body), err)
		return token, fmt.Errorf("GetToken cant unmarshl body of YandexDirectAPI response: %s, error: %v ", body, err)
	}
	return token, nil
}

func (c *yclient) CollectAccountandAddtoBD() (info CreateInfo, err error) {
	resultCamps, err := c.CollectCampaigns()
	if err != nil {
		//if that error occurs, this means, that the new account that user trying
		// to add is agency, so we must collect agency's clients and then
		// add new account in DB for each of agency client
		if strings.Contains(err.Error(), "53") {
			var client goyad.Client
			client.Token = c.Token
			client.Login = c.Login
			client.ApiUrl = c.ApiUrl
			_, err := c.addYandexAgencyAccounts()
			if err != nil {
				logrus.Errorf("CollectAccountandAddtoBD addYandexAgencyAccounts(%s, %s, %s) error: %v", client.Login, client.Token.GetToken(), c, err)
				return info, fmt.Errorf("cant recieve campaings from Yandex.Direct with parametrs %s %s error: ", client.Login, c, err)

			}
			logrus.Info("addYandexAgencyAccounts SUCCESS")
			return info, nil
		} else {
			logrus.Errorln("CollectAccountandAddtoBD unknow error: ", err)
			return info, fmt.Errorf("cant recieve campaings from Yandex.Direct inside CollectAccountandAddtoBD function, error: %v", err)
		}
		return
	}
	resultClientInfo, err := c.CollectClientInfo()
	if err != nil {
		logrus.Println("GetYandexAccessToken GetYandexToken error: ", err)
		return info, fmt.Errorf("cant collectClientInfo from yandex with login: %s, error: %v", c.Login, err)
	}
	logrus.Infof("Information about Client: %+v", resultClientInfo)
	a := model.NewAccount2(
		c.Creator,
		"Яндекс Директ",
		resultClientInfo.Clients[0].Login,
		resultClientInfo.Clients[0].Representatives[0].Email,
	)
	a.Role = "client"
	a.Status = "active"
	a.Owners = append([]string{}, c.Creator)
	a.AccountType = resultClientInfo.Clients[0].Type
	a.AuthToken = c.Token.GetToken()
	a.AppID = Config.YandexDirectAppID
	a.AppSecret = Config.YandexDirectAppSecret
	a.CampaignsInfo = model.AdaptYandexCampaings(resultCamps)
	a.CreatedAt = time.Now()
	err = c.AddAccToDB(*a)
	if err != nil {
		logrus.Errorf("cant a.Update() to DB %v \n error: %v", a.Accountlogin, err)
		return info, fmt.Errorf("cant add account to DB %v \n error: %v", a.Accountlogin, err)
	}
	info.Status = "[Success] Account Added"
	return info, nil
}

func (c *yclient) NewYandexAdder(cc YandexAccountAdder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlValues, err := cc.ParseURL(r.URL)
		if err != nil {
			logrus.Error("Request from YandexOuath received without values")
			http.Error(w, fmt.Sprintf("Request from Yandex received without values %s", urlValues), http.StatusBadRequest)
			return
		}

		code := urlValues["code"]
		token, err := cc.GetToken(code)
		if err != nil {
			logrus.Println("AddYandexAccount GetToken error: ", err)
			http.Error(w, fmt.Sprintf("cant get auth token from Yandex with code: %s, error: %v", code, err), http.StatusBadRequest)
			return
		}
		creator := r.Context().Value("username").(string)
		if creator == "" {
			logrus.Errorf("AddYandexAccount r.Context().Value(username) is empty: ", creator)
			http.Error(w, fmt.Sprintf("Can't identify username inside request context: %s", creator), http.StatusBadRequest)
			return
		}
		c.Creator = creator
		c.Login = urlValues["accountlogin"]
		c.Token = goyad.Token{Value: token.AccessToken}
		_, err = cc.CollectAccountandAddtoBD()
		if err != nil {
			logrus.Errorf("AddYandexAccount r.Context().Value(username) is empty: ", creator)
			http.Error(w, fmt.Sprintf("Can't identify username inside request context: %s", creator), http.StatusBadRequest)
			return
		}
	})
}
func (c *yclient) addYandexAgencyAccounts() (info CreateInfo, err error) {

	var YandexConnectionsLimit = 5
	chAC := make(chan gc.ClientGetItem, 4) // channel's buffer is the number of simultaneous gorouitenes
	var wg sync.WaitGroup
	var client goyad.Client
	client.Token = c.Token
	client.Login = c.Login
	client.ApiUrl = c.ApiUrl
	resultA, err := collectAgencyClients(client)
	if err != nil {
		logrus.Errorln("collectCampaingsfromAgency  error: ", err)
		return info, err
	}
	for i := 0; i < YandexConnectionsLimit; i++ {
		wg.Add(1)
		go func() {
			for {
				ac, ok := <-chAC
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				for _, agencyClient := range ac.Representatives {
					//collecting ads campaigns from Yandex for every agency client
					ci := goyad.NewClient()
					ci.Login = agencyClient.Login
					ci.Token = goyad.Token{Value: client.Token.GetToken()}
					result, err := collectCampaings(ci)
					if err != nil {
						logrus.Errorf("cant collect campaings with parameters: collectCampaings(%s, %s) error: %v", ci.Login, ci.Token.GetToken(), err)
					}

					//creating new account for each of agency clients
					a := model.NewAccount2(
						c.Creator,
						"Яндекс Директ",
						agencyClient.Login,
						agencyClient.Email,
					)
					a.Role = "client"
					a.Status = "active"
					a.Owners = append([]string{}, c.Creator)
					a.AccountType = ac.Type
					a.AuthToken = client.Token.GetToken()
					a.AppID = Config.YandexDirectAppID
					a.AppSecret = Config.YandexDirectAppSecret
					a.CampaignsInfo = model.AdaptYandexCampaings(result)
					a.CreatedAt = time.Now()
					err = c.AddAccToDB(*a)
					if err != nil {
						logrus.Errorf("addYandexAgencyAccounts a.Update(%s) error: %v", agencyClient.Login, err)
						return
					}
				}
			}
		}()
	}
	var agencyClients []string
	for _, ac := range resultA.Clients {
		logrus.Infof("for resultA.Clients ac: %v", ac)
		for _, info := range ac.Representatives {
			agencyClients = append(agencyClients, info.Login)
		}
		chAC <- ac // add client to the queue
	}
	a := model.NewAccount2(
		c.Creator,
		"Яндекс Директ",
		client.Login,
		client.Login+"@yandex.ru",
	)
	a.Creator = c.Creator
	a.Source = "Яндекс Директ"
	a.Accountlogin = client.Login
	a.Email = client.Login + "@yandex.ru"
	a.Role = "agency"
	a.Status = "active"
	a.AgencyClients = agencyClients
	a.Owners = []string{c.Creator}
	a.AuthToken = client.Token.GetToken()
	a.AppID = Config.YandexDirectAppID
	a.AppSecret = Config.YandexDirectAppSecret
	a.CreatedAt = time.Now()
	err = c.AddAccToDB(*a)
	if err != nil {
		logrus.Errorf("cant add account to DB %v \n error: %v", client.Login, err)
		return info, err
	}
	close(chAC) // This tells the goroutines there's nothing else to do
	wg.Wait()   // Wait for the threads to finish
	return info, nil
}
func (c *yclient) CollectClientInfo() (res clients.GetResponse, err error) {
	clientInfo := clients.GetRequest{
		FieldNames: []clients.ClientFieldEnum{
			"ClientId",
			"ClientInfo",
			"CountryId",
			"CreatedAt",
			"Login",
			"Representatives",
			"Type",
		},
	}
	var yc goyad.Client
	yc.Token = c.Token
	yc.Login = c.Login
	yc.ApiUrl = c.ApiUrl
	service2 := clients.New(&yc)
	result, err := service2.Get(clientInfo)
	if err != nil {
		return res, fmt.Errorf("CollectClientInfo service.Get error %v", err)
	}
	return result, nil
}
func (c *yclient) CollectCampaigns() (res campaigns.GetResponse, err error) {
	campaingsInfo := campaigns.GetRequest{
		FieldNames: []campaigns.CampaignFieldEnum{
			"ClientInfo",
			"Id",
			"Name",
			"RepresentedBy",
			"Statistics",
			"Status",
			"Type",
		},
	}
	var yc goyad.Client
	yc.Token = c.Token
	yc.Login = c.Login
	yc.ApiUrl = c.ApiUrl
	service := campaigns.New(&yc)
	result, err := service.Get(campaingsInfo)
	if err != nil {
		return res, fmt.Errorf("collectCampaings service.Get error %v", err)
	}
	return result, nil
}
func (c *yclient) CollectAgencyClients() (res agencyclients.GetResponse, err error) {
	clientInfo := agencyclients.GetRequest{
		FieldNames: []agencyclients.AgencyClientFieldEnum{
			"AccountQuality", "ClientId", "ClientInfo", "Login", "Phone", "Representatives", "Restrictions", "Type",
		},
	}
	var yc goyad.Client
	yc.Token = c.Token
	yc.Login = c.Login
	yc.ApiUrl = c.ApiUrl
	service2 := agencyclients.New(&yc)
	result, err := service2.Get(clientInfo)
	if err != nil {
		return res, fmt.Errorf("collectAgencyClients service.Get error %v", err)
	}
	return result, nil
}
func (c *yclient) AddAccToDB(a model.Account2) error {

	return a.Update()
}
