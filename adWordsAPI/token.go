package adWordsAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bytes"

	"github.com/sirupsen/logrus"
)

//type AdWordsToken struct {
//	AccessToken string `json:"access_token"`
//	ExpiresIn   int    `json:"expires_in"`
//	UserID      int    `json:"user_id"`
//	Email       string `json:"email"`
//}
type AdWordsToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}
type AdWordsError struct {
	Error_           string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e *AdWordsError) Error() string {
	return fmt.Sprintf("Error of VK API : %+v %+v ", e.Error_, e.ErrorDescription)
}

func GetAccessToken(appID, appSecret, redirectURL, authURL, code string) (AdWordsToken, error) {

	var token AdWordsToken
	//https://www.googleapis.com/oauth2/v4/token
	//code=4/P7q7W91a-oMsCeLvIaQm6bTrgtp7&
	//	client_id=your_client_id&
	//	client_secret=your_client_secret&
	//	redirect_uri=https://oauth2.example.com/code&
	//grant_type=authorization_code
	adWordsUrl := authURL +
		"?client_id=" + appID +
		"&client_secret=" + appSecret +
		"&code=" + code +
		"&redirect_uri=" + redirectURL +
		"&grant_type=authorization_code"
	client := &http.Client{}
	r, err := http.NewRequest("POST", adWordsUrl, nil)
	if err != nil {
		logrus.Errorln("GetAccessToken http.NewRequest error: ", err)
		return token, fmt.Errorf("VkAuthorize http.NewRequest error: %s", err)
	}
	logrus.Infof("GetAccessToken AdWORDS Request :\n %+v", r)
	resp, err := client.Do(r)
	if err != nil {
		logrus.Errorln("GetAccessToken client.Do(r) error: ", err)
		return token, fmt.Errorf("VkAuthorize client.Do(r) error: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln("GetAccessToken ioutil.ReadAll(resp.Body) error", err)
		return token, fmt.Errorf("VkAuthorize ioutil.ReadAll(resp.Body) error: %s", err)
	}
	if bytes.Compare(body[:11], []byte("{\"error\":")) == 0 {
		error := AdWordsError{}
		if err := json.Unmarshal(body, &error); err != nil {
			logrus.Errorf("ERROR OF VKAPIError STRUCT inside VKAPI value: %s err:%+v", string(body), err)
			return token, err
		} else {
			return token, &error
		}
	}
	logrus.Infof("GetAccessToken AdWORDS RESPONSE :\n %+v", r)
	logrus.Infof("GetAccessToken AdWORDS RESPONSE body: \n %+v", string(body))
	err = json.Unmarshal(body, &token)
	if err != nil {
		logrus.Error("GetAccessToken json.Unmarshal(body, &token) error: ", err)
		return token, fmt.Errorf("GetAccessToken json.Unmarshal(body, &token) error: %s", err)
	}
	return token, nil
}
