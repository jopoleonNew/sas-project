package adWords

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/adWordsAPI"
)

func AddAdWordsAccount(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("AddAdWordsAccount request : \n %+v ", r)
	query := r.URL.Query()
	if query["code"] == nil || len(query["code"]) == 0 {
		logrus.Error(" Auth Request from AdWords received without code.")
		http.Error(w, fmt.Sprintf(" Auth Request from Vkontakte received without code. %s", query), http.StatusBadRequest)
		return
	}
	code := query["code"]
	logrus.Infof("AddAdWordsAccount r.Query: %+v", query)
	logrus.Infof("AddAdWordsAccount code: %+v", code)
	adtoken, err := adWordsAPI.GetAccessToken(Config.AdWordsAppID, Config.AdWordsAppSecret, Config.AdWordsRedirectURL, "https://www.googleapis.com/oauth2/v4/token", code[0])
	if err != nil {
		logrus.Errorf("AddAdWordsAccount adWordsAPI.GetAccessToken error: %v", err)
		return
	}
	creator := r.Context().Value("username").(string)
	if creator == "" {
		logrus.Errorf("AddAdWordsAccount r.Context().Value(username) is empty: ", creator)
		http.Error(w, fmt.Sprintf("Can't identify username inside AddVKAccount request context: %s", creator), http.StatusBadRequest)
		return
	}
	logrus.Infof("Inside AddAdWordsAccount adWordsAPI.GetAccessToken result:::::: %+v", adtoken)
	return
	//response, err := vk.Request(vktoken.AccessToken, "ads.getAccounts", nil)
	//if err != nil {
	//	logrus.Errorf("AddAdWordsAccount vk.Request error: %v", err)
	//	return
	//}
}
