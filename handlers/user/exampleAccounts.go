package user

import (
	"fmt"
	"net/http"

	"strconv"
	"time"

	"math/rand"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"

	"gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
)

var WA = []string{"nut", "loaf", "haircut", "measure", "cow", "oranges", "event", "twist", "observation", "nose", "tongue", "silk", "aunt", "wool", "hose"}

func AddExmapleAccounts(w http.ResponseWriter, r *http.Request) {

	creator := r.Context().Value("username").(string)
	if creator == "" {
		logrus.Errorf("AddExmapleAccounts r.Context().Value(username) is empty: ", creator)
		http.Error(w, fmt.Sprintf("Can't identify username inside AddExmapleAccounts request context: %s", creator), http.StatusBadRequest)
		return
	}
	for i := 0; i < 3+rand.Intn(5); i++ {
		err := AddVkExampleAccs(creator, "example")
		if err != nil {
			logrus.Errorf("AddVkExampleAccs for creator %v error: %v", creator, err)
			http.Error(w, fmt.Sprintf("can't add VK account %v, \n error: %+v:", err), http.StatusInternalServerError)
			return
		}
	}
	logrus.Infof("\n\n __New Account from Vk added successfully!")
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
	return
}

func DeleteExampleAccs(w http.ResponseWriter, r *http.Request) {
	creator := r.Context().Value("username").(string)
	if creator == "" {
		logrus.Errorf("DeleteExampleAccs r.Context().Value(username) is empty: ", creator)
		http.Error(w, fmt.Sprintf("Can't identify username inside DeleteExampleAccs request context: %s", creator), http.StatusBadRequest)
		return
	}
	a := model.NewAccount2("", "", "", "")
	err := a.DeleteExampleAccs("vk")
	if err != nil {
		logrus.Errorf("DeleteExampleAccs a.DeleteExampleAccs(vk) error: %v", err)
		http.Error(w, fmt.Sprintf("DeleteExampleAccs a.DeleteExampleAccs(yandex) error: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	return
}

func AddVkExampleAccs(creator, email string) error {
	Ncamp := 3 + rand.Intn(3)
	a := model.NewAccount2(creator, "Vkontakte", strconv.Itoa(1000+rand.Intn(1000)), email)
	cmpAr := vkontakteAPI.AdsCampaigns{
		Response: make([]struct {
			ID         int    `json:"id"`
			Type       string `json:"type"`
			Name       string `json:"name"`
			Status     int    `json:"status"`
			DayLimit   string `json:"day_limit"`
			AllLimit   string `json:"all_limit"`
			StartTime  string `json:"start_time"`
			StopTime   string `json:"stop_time"`
			CreateTime string `json:"create_time"`
			UpdateTime string `json:"update_time"`
		}, Ncamp+2),
	}
	for i := 0; i < Ncamp; i++ {
		cmpAr.Response[i].ID = 10000 + rand.Intn(20000)
		cmpAr.Response[i].Type = WA[rand.Intn(len(WA))]
		cmpAr.Response[i].Name = WA[rand.Intn(len(WA))]
		cmpAr.Response[i].Status = rand.Intn(2)
		cmpAr.Response[i].DayLimit = WA[rand.Intn(len(WA))]
		cmpAr.Response[i].AllLimit = WA[rand.Intn(len(WA))]
		cmpAr.Response[i].StartTime = WA[rand.Intn(len(WA))]
		cmpAr.Response[i].StopTime = WA[rand.Intn(len(WA))]
		cmpAr.Response[i].CreateTime = "asdx"
	}

	a.CampaignsInfo = model.AdaptVKCampaings(cmpAr, creator)

	ads := vkontakteAPI.Ads{
		Response: make([]struct {
			ID               string      `json:"id"`
			CampaignID       int         `json:"campaign_id"`
			Name             string      `json:"name"`
			Status           int         `json:"status"`
			Approved         string      `json:"approved"`
			AllLimit         string      `json:"all_limit"`
			CreateTime       string      `json:"create_time"`
			UpdateTime       string      `json:"update_time"`
			AgeRestriction   string      `json:"age_restriction"`
			Category1ID      string      `json:"category1_id"`
			Category2ID      string      `json:"category2_id"`
			CostType         int         `json:"cost_type"`
			AdFormat         int         `json:"ad_format"`
			Cpm              string      `json:"cpm"`
			Cpc              string      `json:"cpc"`
			Video            int         `json:"video"`
			ImpressionsLimit interface{} `json:"impressions_limit"`
			AdPlatform       string      `json:"ad_platform"`
		}, 8),
	}
	for i := 0; i < 2+rand.Intn(3); i++ {
		ads.Response[i].ID = strconv.Itoa(10000 + rand.Intn(20000))
		ads.Response[i].CampaignID = cmpAr.Response[rand.Intn(len(cmpAr.Response))].ID
		ads.Response[i].Name = WA[rand.Intn(len(WA))]
		ads.Response[i].Status = rand.Intn(2)
		ads.Response[i].Approved = WA[rand.Intn(len(WA))]
		ads.Response[i].AllLimit = WA[rand.Intn(len(WA))]
		ads.Response[i].CreateTime = WA[rand.Intn(len(WA))]
		ads.Response[i].UpdateTime = WA[rand.Intn(len(WA))]
		ads.Response[i].AgeRestriction = strconv.Itoa(rand.Intn(90))
		ads.Response[i].Category1ID = "asdx"
		ads.Response[i].Category2ID = "asdx"
		ads.Response[i].CostType = rand.Intn(1)
		ads.Response[i].AdFormat = rand.Intn(1)
		ads.Response[i].Cpm = strconv.Itoa(1000 + rand.Intn(1000))
		ads.Response[i].Cpc = strconv.Itoa(1000 + rand.Intn(1000))
		ads.Response[i].Video = 100 + rand.Intn(100)
		ads.Response[i].ImpressionsLimit = 100 + rand.Intn(100)
		ads.Response[i].AdPlatform = "asdx"
	}
	adaptedAds := model.AdaptVKAds(ads)
	for i, c := range a.CampaignsInfo {
		for _, ad := range adaptedAds {
			if c.ID == ad.CampID {
				a.CampaignsInfo[i].Ads = append(a.CampaignsInfo[i].Ads, ad)
			}
		}
	}

	a.CreatedAt = time.Now()
	a.Role = "client"
	a.Owners = append([]string{}, creator)
	a.AuthToken = "exmaple"
	a.AppID = "Config.VKAppID"
	a.AppSecret = "Config.VKAppSecret"
	err := a.Update()
	if err != nil {
		logrus.Errorf("can't a.Update()  \n error: %v", err)
		return err
	}
	now := time.Now()
	var ctLayout = "2006-01-02"
	amAds := 4 + rand.Intn(4)
	stat := make([]vkontakteAPI.AdStatistic, amAds+1)

	for i := 0; i < amAds; i++ {
		days := 4 - rand.Intn(4)
		randTime := now.AddDate(0, 0, days).Format(ctLayout)
		stat[i].ID = a.CampaignsInfo[rand.Intn(len(a.CampaignsInfo))].ID
		stat[i].Type = "asd"
		stat[i].Stats = make([]struct {
			Day              string      `json:"day"`
			Spent            string      `json:"spent,omitempty"`
			Impressions      interface{} `json:"impressions,omitempty"`
			Clicks           int         `json:"clicks,omitempty"`
			VideoViews       int         `json:"video_views,omitempty"`
			VideoViews_half  int         `json:"video_views_half,omitempty"`
			VideoViews_full  int         `json:"video_views_full,omitempty"`
			VideoClicks_site int         `json:"video_clicks_site,omitempty"`
			JoinRate         int         `json:"join_rate,omitempty"`
		}, days)
		for k := 0; k < days; k++ {
			stat[i].Stats[k].Day = randTime
			stat[i].Stats[k].Spent = strconv.Itoa(rand.Intn(1000))
			stat[i].Stats[k].Impressions = strconv.Itoa(rand.Intn(1000))
			stat[i].Stats[k].Clicks = rand.Intn(1000)
			stat[i].Stats[k].VideoViews = rand.Intn(1000)
			stat[i].Stats[k].VideoViews_half = 0 + rand.Intn(1000)
			stat[i].Stats[k].VideoViews_full = 0 + rand.Intn(1000)
			stat[i].Stats[k].VideoClicks_site = 0 + rand.Intn(1000)
			stat[i].Stats[k].JoinRate = 0 + rand.Intn(1000)
		}
	}
	err = model.SaveVKStatistic(a.Accountlogin, stat)
	if err != nil {
		logrus.Errorf("AddVkExampleAccs model.SaveVKStatistic error: %v", err)

		return err
	}
	return nil
}
