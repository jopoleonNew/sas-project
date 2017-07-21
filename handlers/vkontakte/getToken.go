package vkontakte

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"fmt"

	"strconv"

	"time"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
	vk "gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
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

func GetVKToken(w http.ResponseWriter, r *http.Request) {

	//logrus.Infof("GetVKToken incoming request: %+v", r)
	//logrus.Infof("GetVKToken incoming URL.Fragment: %+v", r.URL.Fragment)
	//log.Println("GetYandexAccessToken income URL query: ", r.URL.Query())
	//API := vk.Api{}
	//API.Request()
	query := r.URL.Query()
	if query["code"] == nil || len(query["code"]) == 0 {
		logrus.Warn("Request from Vkontakte received without code. Making ")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("VkAuthorize ioutil.ReadAll(resp.Body) error", err)
			w.Write([]byte("VkAuthorize ioutil.ReadAll(resp.Body) error" + err.Error()))
			return
		}
		var token vk.VKtoken

		err = json.Unmarshal(body, &token)
		if err != nil {
			logrus.Warn("VkAuthorize bad request repsonse body, trying to unmarshal err ", err)
			var vkerr vk.VKtokenErr
			err = json.Unmarshal(body, &vkerr)
			if err != nil {
				logrus.Errorf("response VkAuthorize VKtokenErr json.Unmarshal error: %+v: \n Indefined body: %+v", err, string(body))
				http.Error(w, fmt.Sprintf("response VkAuthorize VKtokenErr json.Unmarshal error: %+v: \n Indefined body: %+v", err, string(body)), http.StatusBadRequest)
				return
			}
			w.Write([]byte("YandexDirectAPI error: " + vkerr.Error + " " + vkerr.ErrorDesсription))
			return
		}
		logrus.Info("Inside VKauthorize vk.VkAccessToken result:::::: ", token)
	}
	code := query["code"]
	vktoken, err := vk.VkAccessToken(Config.VKAppID, Config.VKAppSecret, Config.VKRedirectURL, code[0])
	if err != nil {
		logrus.Println("GetVKToken vk.VkAccessToken error: ", err)
		return
	}
	creator := r.Context().Value("username").(string)
	if creator == "" {
		logrus.Errorf("GetVKToken r.Context().Value(username) is empty: ", creator)
		http.Error(w, fmt.Sprintf("Can't identify username inside GetVKToken request context: %s", creator), http.StatusBadRequest)
		return
	}
	logrus.Info("Inside VKauthorize vk.VkAccessToken result:::::: ", vktoken)
	response, err := vk.Request(vktoken.AccessToken, "ads.getAccounts", nil)
	if err != nil {
		logrus.Println("VKauthorize vk.Request error: ", err)
		return
	}
	logrus.Info("vk.Request ads.getAccounts result: ", string(response))
	var accounts vk.AdsAccounts
	if err := json.Unmarshal(response, &accounts); err != nil {
		logrus.Errorf("can't unmarshal VK response from ads.getAccounts, error: &v", err)
		http.Error(w, fmt.Sprintf("can't unmarshal VK response from ads.getAccounts error: %+v:", err), http.StatusBadRequest)
		return
	}
	// creating NewUser to get info about creator
	// from DB and use it's email to create account with creator's email
	user := model.NewUser()
	user.Username = creator
	userInfo, err := user.GetInfo()
	if err != nil {
		logrus.Errorf("GetVKToken user.GetInfo(%s) error: %v", creator, err)
		userInfo.Email = vktoken.Email
	}

	for _, acc := range accounts.Response {
		//creating new account for adding VK account in DB

		// Depending on what type of account is it, collecting campaings from Vk API:
		// basic account
		if acc.AccountType == "general" {
			a := model.NewAccount2(creator, "Vkontakte", strconv.Itoa(acc.AccountID), userInfo.Email)
			p := make(map[string]string)
			p["account_id"] = strconv.Itoa(acc.AccountID)
			resp, err := vk.Request(vktoken.AccessToken, "ads.getCampaigns", p)
			if err != nil {
				logrus.Println("VKauthorize vk.Request error: ", err)
				return
			}
			logrus.Errorf("VK response from ads.getCampaigns, error: &+v", string(resp))
			var camps vk.AdsCampaigns
			if err := json.Unmarshal(resp, &camps); err != nil {
				logrus.Errorf("can't unmarshal VK response from ads.getCampaigns, error: &v", err)
				http.Error(w, fmt.Sprintf("can't unmarshal VK response from ads.getCampaigns error: %+v:", err), http.StatusBadRequest)
				return
			}
			a.CampaignsInfo = model.AdaptVKCampaings(camps, strconv.Itoa(acc.AccountID))
			a.CreatedAt = time.Now()
			if acc.AccountType == "general" {
				a.Role = "client"
			}
			if acc.AccountType == "agency" {
				a.Role = "agency"
			}
			a.Owners = append([]string{}, creator)
			a.AuthToken = vktoken.AccessToken
			a.AppID = Config.VKAppID
			a.AppSecret = Config.VKAppSecret
			err = a.AdvanceUpdate()
			if err != nil {
				logrus.Println("a.AdvanceUpdate() for AccountType = general error: ", err)
				return
			}
			return
		}
		// agency account
		if acc.AccountType == "agency" {
			p := make(map[string]string)
			p["account_id"] = strconv.Itoa(acc.AccountID)
			resp, err := vk.Request(vktoken.AccessToken, "ads.getClients", p)
			if err != nil {
				logrus.Println("VKauthorize vk.Request error: ", err)
				return
			}
			var clients vk.AdsClients
			if err := json.Unmarshal(resp, &clients); err != nil {
				logrus.Errorf("can't unmarshal VK response from ads.getClients, error: &v", err)
				http.Error(w, fmt.Sprintf("can't unmarshal VK response from ads.getClients error: %+v:", err), http.StatusBadRequest)
				return
			}
			for _, client := range clients.Response {
				p := make(map[string]string)
				p["account_id"] = strconv.Itoa(acc.AccountID)
				p["client_id"] = strconv.Itoa(client.ID)
				campaigns, err := vk.Request(vktoken.AccessToken, "ads.getCampaigns", nil)
				if err != nil {
					logrus.Println("VKauthorize vk.Request error: ", err)
					return
				}
				var camps vk.AdsCampaigns
				if err := json.Unmarshal(campaigns, &camps); err != nil {
					logrus.Errorf("can't unmarshal VK response from ads.getCampaigns, error: &v", err)
					http.Error(w, fmt.Sprintf("can't unmarshal VK response from ads.getCampaigns error: %+v:", err), http.StatusBadRequest)
					return
				}
				a := model.NewAccount2(creator, "Vkontakte", strconv.Itoa(acc.AccountID), userInfo.Email)
				a.CampaignsInfo = model.AdaptVKCampaings(camps, client.Name)
				a.CreatedAt = time.Now()
				if acc.AccountType == "general" {
					a.Role = "client"
				}
				if acc.AccountType == "agency" {
					a.Role = "agency"
				}
				a.Owners = append([]string{}, client.Name, strconv.Itoa(acc.AccountID))
				a.AuthToken = vktoken.AccessToken
				a.AppID = Config.VKAppID
				a.AppSecret = Config.VKAppSecret
				err = a.AdvanceUpdate()
				if err != nil {
					logrus.Println("a.AdvanceUpdate() AccountType = agency error: ", err)
					return
				}
				return
			}
			a := model.NewAccount2(creator, "Vkontakte", strconv.Itoa(acc.AccountID), userInfo.Email)
			a.CampaignsInfo = model.AdaptVKCampaings(vk.AdsCampaigns{}, strconv.Itoa(acc.AccountID))
			a.CreatedAt = time.Now()
			if acc.AccountType == "general" {
				a.Role = "client"
			}
			if acc.AccountType == "agency" {
				a.Role = "agency"
			}
			a.Owners = append([]string{}, creator)
			a.AuthToken = vktoken.AccessToken
			a.AppID = Config.VKAppID
			a.AppSecret = Config.VKAppSecret
			err = a.AdvanceUpdate()
			if err != nil {
				logrus.Println("a.AdvanceUpdate() AccountType = agency error: ", err)
				return
			}
			return
		}
	}
	return
	//{"response":[{"account_id":1604082246,"account_type":"general","account_status":1,"access_role":"reports"}]}

}
