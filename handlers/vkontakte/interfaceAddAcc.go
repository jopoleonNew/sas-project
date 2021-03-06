package vkontakte

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"net/url"

	"gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
)

type VKAccount struct {
	RedirectURL string
	ApiURL      string
	AuthURL     string
	AppID       string
	AppSecret   string
	Email       string
	params      map[string]string
}
type vkToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	Email     string `json:"email"`
}

func (t *vkToken) GetToken() string  { return t.Token }
func (t *vkToken) GetExpiresIn() int { return t.ExpiresIn }
func (t *vkToken) GetEmail() string  { return t.Email }

func (vka *VKAccount) GetAccessToken(code string) (vkToken, error) {
	vt, err := vkontakteAPI.GetAccessToken(Config.VKAppID, Config.VKAppSecret, Config.VKRedirectURL, vka.AuthURL, code)
	if err != nil {
		logrus.Errorln("VkAuthorize client.Do(r) error: ", err)
		return vkToken{}, fmt.Errorf("VkAuthorize client.Do(r) error: %s", err)
	}
	return vkToken{
		Token:     vt.AccessToken,
		ExpiresIn: vt.ExpiresIn,
		Email:     vt.Email,
	}, nil

}
func (vka *VKAccount) CollectCampaings(token string) ([]byte, error) {
	return []byte{}, nil
}
func (vka *VKAccount) ParseURL(url *url.URL) (map[string]string, error) {
	var p map[string]string
	query := url.Query()
	s := query["code"]
	if s == nil || len(s) == 0 {
		return p, fmt.Errorf("Request received without code %s", query)
	}
	al := query["state"]
	if al == nil || len(al) == 0 {
		return p, fmt.Errorf("Request received without state %s", query)
	}
	p["code"] = s[0]
	p["accountlogin"] = al[0]
	return p, nil
}
