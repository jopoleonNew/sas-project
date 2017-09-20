package vkontakte

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"net/url"

	"gogs.itcloud.pro/SAS-project/sas/model"
)

type AccountAdder interface {
	ParseURL(url *url.URL) (map[string]string, error)
	//SetParams() (params map[string]string)
	GetAccessToken(code string) (TokenInterface, error)
	InfoCollecter
	UpdateInDB(a *model.Account2) error
}

type InfoCollecter interface {
	//GetAccountInfo(token string, params map[string]string) ([]byte, error)
	ClientCollector
	AgencyCollector
	CollectCampaings(token string) ([]byte, error)
	CollectAds(token string) ([]byte, error)
}
type ClientCollector interface {
	CollectClientInfo(token string) (model.Account2, error)
}
type AgencyCollector interface {
	CollectAgencyInfo(token string) (model.Account2, error)
}

//type AccountInfoCollector interface {
//}

type TokenInterface interface {
	GetToken() string
	GetExpiresIn() int
	GetEmail() string
	//Yandex
	//TokenType    string `json:"token_type"`
	//AccessToken  string `json:"access_token"`
	//ExpiresIn    int    `json:"expires_in"`
	//RefreshToken string `json:"refresh_token"`

	//VK
	//AccessToken string `json:"access_token"`
	//ExpiresIn   int    `json:"expires_in"`
	//UserID      int    `json:"user_id"`
	//Email       string `json:"email"`
}
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

//func (vka *VKAccount) GetAccessToken(appID, appSecret, code string) ([]byte, error) {
//	vkontakteAPI.GetAccessToken(cfg.VKAppID, cfg.VKAppSecret, cfg.VKRedirectURL, vka.AuthURL, code[0])
//
//}
//func (vka *VKAccount) CollectCampaings(token string, params map[string]string) ([]byte, error) {}

func NewAdderHandler(cc AccountAdder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		query := r.URL.Query()
		logrus.Info("AddYandexAccount income URL query: ", r.URL.Query())
		urlParams, err := cc.ParseURL(r.URL)
		code, ok := urlParams["code"]
		if !ok {
			logrus.Errorf("Request received without code %s", query)
			http.Error(w, fmt.Sprintf("Request received without code %s", query), http.StatusBadRequest)
			return
		}
		accountlogin, ok := urlParams["accountlogin"]
		if !ok {
			logrus.Errorf("Request received without accountlogin %s", query)
			http.Error(w, fmt.Sprintf("Request received without accountlogin %s", query), http.StatusBadRequest)
			return
		}
		logrus.Info("AccountLogin recieved: ", accountlogin)
		token, err := cc.GetAccessToken(code)
		if err != nil {
			http.Error(w, fmt.Sprintf("GetAccessToken error: %v", err), http.StatusBadRequest)
			return
		}
		creator := r.Context().Value("username").(string)
		if creator == "" {
			logrus.Errorf("NewAdderHandler r.Context().Value(username) is empty: ", creator)
			http.Error(w, fmt.Sprintf("Can't identify username inside AddVKAccount request context: %s", creator), http.StatusBadRequest)
			return
		}
		cc.CollectClientInfo(token.GetToken())
	})
}
