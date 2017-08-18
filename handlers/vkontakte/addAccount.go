package vkontakte

import (
	"encoding/json"
	"net/http"

	"fmt"

	"strconv"

	"time"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
	vk "gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
)

func AddVKAccount(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query["code"] == nil || len(query["code"]) == 0 {
		logrus.Error(" Auth Request from Vkontakte received without code.")
		http.Error(w, fmt.Sprintf(" Auth Request from Vkontakte received without code. %s", query), http.StatusBadRequest)
		return
	}
	code := query["code"]
	vktoken, err := vk.GetVKAccessToken(Config.VKAppID, Config.VKAppSecret, Config.VKRedirectURL, "https://oauth.vk.com/access_token", code[0])
	if err != nil {
		logrus.Errorf("AddVKAccount vk.GetVKAccessToken error: %v", err)
		return
	}
	creator := r.Context().Value("username").(string)
	if creator == "" {
		logrus.Errorf("AddVKAccount r.Context().Value(username) is empty: ", creator)
		http.Error(w, fmt.Sprintf("Can't identify username inside AddVKAccount request context: %s", creator), http.StatusBadRequest)
		return
	}
	logrus.Info("Inside VKauthorize vk.GetVKAccessToken result:::::: ", vktoken)
	response, err := vk.Request(vktoken.AccessToken, "ads.getAccounts", nil)
	if err != nil {
		logrus.Errorf("VKauthorize vk.Request error: %v", err)
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
		logrus.Errorf("AddVKAccount user.GetInfo(%s) error: %v", creator, err)
		userInfo.Email = vktoken.Email
	}

	for _, acc := range accounts.Response {
		// Depending on what type of account is it, collecting campaings from Vk API:
		// basic account
		if acc.AccountType == "general" {
			err := addGeneralAccount(acc, vktoken.AccessToken, creator, userInfo.Email)
			if err != nil {
				logrus.Errorf("addGeneralAccount error: %v", creator, err)
				http.Error(w, fmt.Sprintf("can't add VK account %v, \n error: %+v:", acc, err), http.StatusBadRequest)
				return
			}
		}
		// agency account
		if acc.AccountType == "agency" {
			logrus.Errorf("\n\n Adding Agency account is not implemented yet.")
		}
	}
	logrus.Infof("\n\n __New Account from Vk added successfully!")
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
	return

}

type VKCollector interface {
	collectCampaigns(token string, params map[string]string) ([]byte, error)
	collectAds(token string, params map[string]string) ([]byte, error)
}

func CollectAccountsandAddToDB() {

}
func addAgencyAccount(acc vk.AccountList, token, creator, email string) error {
	p := make(map[string]string)
	p["account_id"] = strconv.Itoa(acc.AccountID)
	//getting the list of agency clients
	resp, err := vk.Request(token, "ads.getClients", p)
	if err != nil {
		logrus.Errorf("VKauthorize vk.Request error: %v", err)
		return err
	}
	var clients vk.AdsClients
	if err := json.Unmarshal(resp, &clients); err != nil {
		logrus.Errorf("can't unmarshal VK response from ads.getClients, error: %v", err)
		return err
	}
	var agencyClients []string
	for _, client := range clients.Response {
		agencyClients = append(agencyClients, strconv.Itoa(client.ID))
		p := make(map[string]string)
		p["account_id"] = strconv.Itoa(acc.AccountID)
		p["client_id"] = strconv.Itoa(client.ID)
		camps, err := collectCampaigns(token, p)
		if err != nil {
			logrus.Errorf("VKauthorize vk.Request error: %v", err)
			return fmt.Errorf("can't get ads.getCampaigns from VK, error: %v", err)
		}

		a := model.NewAccount2(creator, "Vkontakte", strconv.Itoa(acc.AccountID), email)
		a.CampaignsInfo = model.AdaptVKCampaings(camps, client.Name)
		a.CreatedAt = time.Now()
		if acc.AccountType == "general" {
			a.Role = "client"
		}
		if acc.AccountType == "agency" {
			a.Role = "agency"
		}
		a.Owners = append([]string{}, client.Name, strconv.Itoa(acc.AccountID))
		a.AuthToken = token
		a.AppID = Config.VKAppID
		a.AppSecret = Config.VKAppSecret
		err = a.Update()
		if err != nil {
			logrus.Errorf("a.Update() AccountType = agency error: ", err)
			return fmt.Errorf("can't a.Update() account: %v, \n error: %v", a, err)
		}
	}
	a := model.NewAccount2(creator, "Vkontakte", strconv.Itoa(acc.AccountID), email)
	a.CampaignsInfo = model.AdaptVKCampaings(vk.AdsCampaigns{}, strconv.Itoa(acc.AccountID))
	a.CreatedAt = time.Now()
	if acc.AccountType == "general" {
		a.Role = "client"
	}
	if acc.AccountType == "agency" {
		a.Role = "agency"
	}
	a.AgencyClients = agencyClients
	a.Owners = append([]string{}, creator)
	a.AuthToken = token
	a.AppID = Config.VKAppID
	a.AppSecret = Config.VKAppSecret
	err = a.Update()
	if err != nil {
		logrus.Errorf("a.Update() AccountType = agency error: ", err)
		return fmt.Errorf("can't a.Update() account: %v, \n error: %v", a, err)
	}
	return nil
}

//addGeneralAccount creates new account for adding VK general account in DB
func addGeneralAccount(acc vk.AccountList, token, creator, email string) error {
	a := model.NewAccount2(creator, "Vkontakte", strconv.Itoa(acc.AccountID), email)
	p := make(map[string]string)
	p["account_id"] = strconv.Itoa(acc.AccountID)
	camps, err := collectCampaigns(token, p)
	if err != nil {
		logrus.Errorf("can't collectCampaigns for account %v, \n error: %v", acc, err)
		return err
	}

	a.CampaignsInfo = model.AdaptVKCampaings(camps, strconv.Itoa(acc.AccountID))

	for i, c := range a.CampaignsInfo {
		time.Sleep(700 * time.Millisecond)
		p["campaign_id"] = "{\"" + strconv.Itoa(c.ID) + "\"}"
		ads, err := collectAds(token, p)
		if err != nil {
			logrus.Errorf("can't collectAds for account %v, \n error: %v", acc, err)
			return err
		}
		logrus.Infof("Collected Ads for campaing %v , : \n %s", c, ads)
		c.Ads = model.AdaptVKAds(ads)
		a.CampaignsInfo[i] = c
	}

	a.CreatedAt = time.Now()
	if acc.AccountType == "general" {
		a.Role = "client"
	}
	if acc.AccountType == "agency" {
		a.Role = "agency"
	}
	a.Owners = append([]string{}, creator)
	a.AuthToken = token
	a.AppID = Config.VKAppID
	a.AppSecret = Config.VKAppSecret
	err = a.Update()
	if err != nil {
		logrus.Errorf("can't a.Update() for %a, \n error: %v", acc, err)
		return err
	}
	return nil
}

//collectCampaigns collects advertisement campaigns from VK API
func collectCampaigns(token string, params map[string]string) (vk.AdsCampaigns, error) {

	var camps vk.AdsCampaigns
	resp, err := vk.Request(token, "ads.getCampaigns", params)
	if err != nil {
		logrus.Errorf("VKauthorize vk.Request error: %v", err)
		return camps, fmt.Errorf("collectCampaigns vk.Request error: %v", err)
	}
	//logrus.Errorf("VK response from ads.getCampaigns, error: &+v", string(resp))

	if err := json.Unmarshal(resp, &camps); err != nil {
		logrus.Errorf("can't unmarshal VK response from ads.getCampaigns, error: &v", err)

		return camps, fmt.Errorf("collectCampaigns json.Unmarshal error: %v", err)
	}
	return camps, nil
}

//account_idидентификатор рекламного кабинета.
//обязательный параметр, целое число
//client_idДоступно и обязательно для рекламных агентств. Идентификатор клиента, у которого запрашиваются рекламные объявления.
//целое число
//include_deletedФлаг, задающий необходимость вывода архивных объявлений.
//0 — выводить только активные объявления;
//1 — выводить все объявления.
//флаг, может принимать значения 1 или 0
//campaign_idsфильтр по рекламным кампаниям.
//Сериализованный JSON-массив, содержащий id кампаний. Если параметр равен null, то будут выводиться рекламные объявления всех кампаний.
//строка
//ad_idsфильтр по рекламным объявлениям.
//Сериализованный JSON-массив, содержащий id объявлений. Если параметр равен null, то будут выводиться все рекламные объявления.
//строка
//limitограничение на количество возвращаемых объявлений. Используется, только если параметр ad_ids равен null, а параметр campaign_ids содержит id только одной кампании.
//целое число
//offsetсмещение. Используется в тех же случаях, что и параметр limit.
//целое число
func collectAds(token string, params map[string]string) (vk.Ads, error) {
	var ads vk.Ads
	resp, err := vk.Request(token, "ads.getAds", params)
	if err != nil {
		logrus.Errorf("collectAds vk.Request error: %v", err)
		return ads, fmt.Errorf("collectAds vk.Request error: %v", err)
	}

	if err := json.Unmarshal(resp, &ads); err != nil {
		logrus.Errorf("can't unmarshal VK response from ads.getAds, error: &v", err)

		return ads, fmt.Errorf("collectAds json.Unmarshal error: %v", err)
	}
	return ads, nil
}
