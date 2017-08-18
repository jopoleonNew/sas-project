package vkontakteAPI

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"fmt"

	"bytes"

	"encoding/json"

	"github.com/sirupsen/logrus"
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
	Method string            `json:"method"`
	Token  string            `json:"token"`
	Values map[string]string `json:"values"`
}

//https://api.vk.com/method/METHOD_NAME?PARAMETERS&access_token=ACCESS_TOKEN&v=V
// Request makes request to VK API with given method name and parameters
// You can see full list of them in official docs https://vk.com/dev/manuals
func Request(token, methodName string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(API_METHOD_URL + methodName)
	if err != nil {
		logrus.Errorf("VkAPI Request url.Parse error: %v", err)
		return nil, fmt.Errorf("url.Parse error: %v", err)
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	q.Set("access_token", token)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		logrus.Errorf("VkAPI Request http.Get(u.String() error: %v", err)
		return nil, fmt.Errorf("http.Get(u.String()) error: %v", err)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("VkAPI Request ioutil.ReadAll(resp.Body) error: %v", err)
		return nil, fmt.Errorf("VkAPI Request ioutil.ReadAll(resp.Body) error: %v", err)
	}
	if bytes.Compare(content[:9], []byte("{\"error\":")) == 0 {
		error := VKAPIError{}
		if err := json.Unmarshal(content, &error); err != nil {
			logrus.Errorf("ERROR OF VKAPIError STRUCT inside VKAPI value: %s err:%+v", string(content), err)
			return content, err
		} else {
			return content, &error
		}
	}

	return content, nil
}

type VKAPIError struct {
	Error_ struct {
		ErrorCode     int    `json:"error_code"`
		ErrorMsg      string `json:"error_msg"`
		RequestParams []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"request_params"`
	} `json:"error"`
}

func (e *VKAPIError) Error() string {
	return fmt.Sprintf("Error of VK API : %+v ", e.Error_)
}

type AdsAccounts struct {
	Response []struct {
		AccountID     int    `json:"account_id"`
		AccountType   string `json:"account_type"`
		AccountStatus int    `json:"account_status"`
		AccessRole    string `json:"access_role"`
	} `json:"response"`
}
type AccountList struct {
	AccountID     int    `json:"account_id"`
	AccountType   string `json:"account_type"`
	AccountStatus int    `json:"account_status"`
	AccessRole    string `json:"access_role"`
}
type AdsCampaigns struct {
	Response []struct {
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
	} `json:"response"`
}
type Ads struct {
	Response []struct {
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
	} `json:"response"`
}

//id — идентификатор клиента;
//name — название клиента;
//day_limit — дневной лимит клиента в рублях;
//all_limit — общий лимит клиента в рублях.
type AdsClients struct {
	Response []struct {
		ID        int    `json:"client_id"`
		Name      string `json:"client_name"`
		day_limit int    `json:"day_limit"`
		all_limit int    `json:"all_limit"`
	} `json:"response"`
}
type CampaignsStat struct {
	Response []struct {
		ID    int    `json:"id"`
		Type  string `json:"type"`
		Stats []struct {
			Overall     int    `json:"overall"`
			DayFrom     string `json:"day_from"`
			DayTo       string `json:"day_to"`
			Spent       string `json:"spent"`
			Impressions int    `json:"impressions"`
			Clicks      int    `json:"clicks"`
		} `json:"stats"`
	} `json:"response"`
}
