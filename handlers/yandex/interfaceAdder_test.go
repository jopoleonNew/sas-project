package yandex

import (
	"net/url"
	"testing"

	"reflect"

	"net/http"
	"net/http/httptest"

	"encoding/json"

	"io/ioutil"

	"fmt"

	"errors"

	"github.com/nk2ge5k/goyad"
	"github.com/nk2ge5k/goyad/agencyclients"
	"github.com/nk2ge5k/goyad/campaigns"
	"github.com/nk2ge5k/goyad/clients"
	"github.com/nk2ge5k/goyad/gc"
	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

//type yclient struct {
//	goyad.Client
//	AuthURL     string
//	RedirectURL string
//	ApiURL      string
//	AppID       string
//	AppSecret   string
//	Creator     string
//}
func TestYclient_ParseURL(t *testing.T) {
	//ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Add("Requestid", "12344")
	//	fmt.Fprintln(w, string(testAccount))
	//}))
	//
	//defer ts.Close()
	var tc yclient
	tc.Token = goyad.Token{Value: "123"}
	tc.ApiURL = "localhost:8080"
	tc.Login = "t"

	tt := []struct {
		query     string
		expeceted map[string]string
	}{
		{"code=123&state=t", map[string]string{"code": "123", "accountlogin": "t"}},
		{"code=321&state=x", map[string]string{"code": "321", "accountlogin": "x"}},
	}
	for _, tcc := range tt {
		var u url.URL
		u.Scheme = "http"
		u.Host = "test"
		u.Path = "parseurl"
		u.RawQuery = tcc.query
		p, err := tc.ParseURL(&u)
		if err != nil {
			t.Fatalf("ParseURL error: %v", err)
		}
		if !reflect.DeepEqual(p, tcc.expeceted) {
			t.Fatalf("got: %s, expceted: %s", p, tcc.expeceted)
		}
	}
	btt := []struct {
		query     string
		expeceted map[string]string
	}{

		{"code=1", map[string]string{"code": "1", "accountlogin": ""}},
		{"state=x", map[string]string{"code": "", "accountlogin": "x"}},
	}
	for _, tcc := range btt {
		var u url.URL
		u.Scheme = "http"
		u.Host = "test"
		u.Path = "parseurl"
		u.RawQuery = tcc.query
		_, err := tc.ParseURL(&u)
		if err == nil {
			t.Fatalf("ParseURL returned no error on %v", tcc.query)
		}
	}
}

func TestYclient_GetToken(t *testing.T) {
	var tcli yclient
	//tcli.Token = goyad.Token{Value: "123"}
	tcli.ApiURL = "/"
	tcli.Login = "t"
	//tt := []struct {
	//	url string
	//	res struct {
	//		token string
	//		err   error
	//	}
	//	expected struct {
	//		msg string
	//		err error
	//	}
	//}{
	//	{url: "/good", res: struct {
	//		token string
	//		err   error
	//	}{token: "123", err: errors.New("no error")}},
	//	{url: "/returnerr", res: struct {
	//		token string
	//		err   error
	//	}{token: "123", err: errors.New("no error")}},
	//	{url: "/codeerr", res: struct {
	//		token string
	//		err   error
	//	}{token: "123", err: errors.New("no error")}},
	//	{url: "/bad", res: struct {
	//		token string
	//		err   error
	//	}{token: "123", err: errors.New("no error")}},
	//}
	tserv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Requestid", "12344")

		switch r.URL.Path {
		case "/good":
			var respWriter yad.YandexTokenbody
			respWriter.TokenType = "123"
			respWriter.AccessToken = "123"
			respWriter.ExpiresIn = 3
			respWriter.RefreshToken = "123"
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error("NewServer cant read incoming request body, error: ", err)
			}
			goodBody := "client_id=&client_secret=&code=11113333&grant_type=authorization_code"
			if string(body) != goodBody {
				t.Errorf("NewServer wrong body recieved, got %s, expected %s ", string(body), goodBody)
			}
			defer r.Body.Close()
			b, err := json.Marshal(respWriter)
			if err != nil {
				t.Fatalf("NewServer unmarshal yandex token struct, error: %v", err)
			}
			w.Write(b)
		case "/returnerr":
			var respWriter yad.YandexV5Error
			respWriter.ErrorDescription = "ab"
			respWriter.ErrorCode = "ab"
			respWriter.ErrorString = "ab"
			respWriter.RequsetID = "ab"
			b, err := json.Marshal(respWriter)
			if err != nil {
				t.Fatalf("NewServer unmarshal yandex token struct, error: %v", err)
			}
			w.Write(b)
		case "/codeerr":
			w.Write([]byte(`{"error_description": "Invalid code", "error": "bad_verification_code"}`))
		case "/bad":
			fmt.Fprintln(w, []byte("cant unmarshal"))

		}
	}))
	defer tserv.Close()

	tcli.AuthURL = tserv.URL + "/good"
	logrus.Info(tcli.ApiURL)
	resp, err := tcli.GetToken("11113333")
	if err != nil {
		t.Fatalf("tcli.GetToken error: %v", err)
	}
	if resp.AccessToken != "123" {
		t.Fatalf("Wrong response token, got %s, expeceted %s", resp.AccessToken, "123")
	}
	tcli.AuthURL = tserv.URL + "/returnerr"
	logrus.Info(tcli.ApiURL)
	_, err = tcli.GetToken("11113333")
	if err == nil {
		t.Fatalf("tcli.GetToken didnt return error %v", err)
	}

	tcli.AuthURL = tserv.URL + "/codeerr"
	logrus.Info(tcli.ApiURL)
	_, err = tcli.GetToken("11113333")
	if err.Error() != "Error from yandexDirect API: Invalid code, bad_verification_code" {
		t.Fatalf("tcli.GetToken return unexpected error, got: %v, execeted: %s", err, "Error from yandexDirect API: Invalid code, bad_verification_code")
	}
	tcli.AuthURL = tserv.URL + "/bad"
	logrus.Info(tcli.ApiURL)
	_, err = tcli.GetToken("11113333")
	if err == nil {
		t.Fatalf("tcli.GetToken didnt return error")
	}
	//logrus.Info(resp)
}
func TestYclient_CollectCampaigns(t *testing.T) {
	testAccount, err := ioutil.ReadFile("testAccount.json")
	if err != nil {
		t.Fatalf("cant read file with test json, error: %v", err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Requestid", "12344")
		fmt.Fprintln(w, string(testAccount))
	}))
	defer ts.Close()
	var tcli yclient
	tcli.ApiUrl = ts.URL + "/"
	tcli.Login = "t"
	tcli.Token = goyad.Token{Value: "123"}
	resp, err := tcli.CollectCampaigns()
	if err != nil {
		t.Fatalf("tcli.CollectCampaigns() error: %v", err)
	}
	expceted := map[string]campaigns.GetResponse{}
	if err := json.Unmarshal(testAccount, &expceted); err != nil {
		t.Fatalf("bad response from collectCampaings, can't unmarshal: %v", err)
	}
	if expectedRes, ok := expceted["result"]; ok {

		var res, expected []byte
		if res, err = json.Marshal(resp.Campaigns); err != nil {
			t.Fatalf("cant marshal : %v", err)
		}
		if expected, err = json.Marshal(expectedRes.Campaigns); err != nil {
			t.Fatalf("cant marshal : %v", err)
		}
		if !reflect.DeepEqual(res, expected) {
			t.Fatalf("result is not equal to expected value")
		}
	} else {
		t.Fatalf("can't find expceted[result] value in result map: ", expectedRes)
	}
}
func TestYclient_CollectAgencyClients(t *testing.T) {
	testAgency := map[string]agencyclients.GetResponse{
		"result": {Clients: []gc.ClientGetItem{

			{
				Login: "qwe",
				Type:  "client",
				Representatives: []gc.Representative{
					{
						Email: "e",
						Login: "e",
						Role:  "e",
					},
				},
			},
		},
		}}

	//testAgency, err := ioutil.ReadFile("testClients.json")
	//if err != nil {
	//	t.Fatalf("cant read file with test json, error: %v", err)
	//}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Requestid", "12344")

		b, err := json.Marshal(&testAgency)
		if err != nil {
			t.Fatalf("json.Marshal error: %v", err)
		}
		logrus.Info("Inside test handler:", string(b))
		fmt.Fprintln(w, b)
	}))
	defer ts.Close()
	var tcli yclient
	tcli.ApiUrl = ts.URL + "/"
	tcli.Login = "t"
	tcli.Token = goyad.Token{Value: "123"}
	resp, err := tcli.CollectAgencyClients()
	if err != nil {
		t.Fatalf("tcli.CollectCampaigns() error: %v", err)
	}
	logrus.Info("Is Equal Agency? : ", reflect.DeepEqual(resp, testAgency))
}
func TestYclient_CollectClientInfo(t *testing.T) {
	testClientInfo := map[string]clients.GetResponse{
		"result": {Clients: []gc.ClientGetItem{
			{
				Login: "qwe",
				Type:  "client",
				Representatives: []gc.Representative{
					{
						Email: "e",
						Login: "e",
						Role:  "e",
					},
				},
			},
		},
		},
	}

	//testClientInfo, err := ioutil.ReadFile("testClients.json")
	//if err != nil {
	//	t.Fatalf("cant read file with test json, error: %v", err)
	//}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Requestid", "12344")

		b, err := json.Marshal(testAgencyClients)
		if err != nil {
			t.Fatalf("json.Marshal error: %v", err)
		}
		logrus.Info("Inside test handler:", string(b))
		fmt.Fprintln(w, b)
	}))
	defer ts.Close()
	var tcli yclient
	tcli.ApiUrl = ts.URL + "/"
	tcli.Login = "t"
	tcli.Token = goyad.Token{Value: "123"}
	resp, err := tcli.CollectClientInfo()
	if err != nil {
		t.Fatalf("tcli.CollectCampaigns() error: %v", err)
	}
	logrus.Info("Is Equal Agency? : ", reflect.DeepEqual(resp, testClientInfo))
}

type FakeAdder struct{}

func (f *FakeAdder) AddAccToDB(a model.Account2) error {
	logrus.Info("Facing adding to DB")
	return nil
}
func (f *FakeAdder) AddYandexAgencyAccounts(db AccountAdder) (info CreateInfo, err error) {
	return info, nil
}
func (f *FakeAdder) CollectAgencyClients() (res agencyclients.GetResponse, err error) {
	res = agencyclients.GetResponse{
		Clients: []gc.ClientGetItem{

			{
				Login: "qwe",
				Type:  "client",
				Representatives: []gc.Representative{
					{
						Email: "e",
						Login: "e",
						Role:  "e",
					},
				},
			},
		},
	}
	return res, nil
}
func (f *FakeAdder) CollectClientInfo() (res clients.GetResponse, err error) {
	return
}
func (f *FakeAdder) CollectCampaigns() (res campaigns.GetResponse, err error) {

	return res, errors.New("53")
}

func TestYclient_CollectAccountandAddtoBD(t *testing.T) {
	var tcli yclient
	//}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Requestid", "12344")
		fmt.Fprintln(w, "{\"result\":{53}}")
	}))
	defer ts.Close()
	//tcli.Token = goyad.Token{Value: "123"}
	tcli.ApiUrl = ts.URL + "/"
	tcli.Login = "t"
	tcli.Token = goyad.Token{Value: "123"}
	fakeDb := FakeAdder{}
	_, err := tcli.CollectAccountandAddtoBD(&fakeDb)
	if err != nil {
		t.Fatalf("tcli.CollectAccountandAddtoBD() error: %v", err)
	}
}

var testAgencyClients = `{
	"result": {
		"ApiObject": {
			"ForceFields": [],
			"NullFields": []
		},
		"GetResponseGeneral": {
			"ApiObject": {
				"ForceFields": [],
				"NullFields": []
			},
			"LimitedBy": "nil"
		},
		"Clients": [
			{
				"ApiObject": {
					"ForceFields": [],
					"NullFields": []
				},
				"ClientBaseItem": {
					"ApiObject": {
						"ForceFields": [],
						"NullFields": []
					},
					"ClientInfo": "Цветы мытищи",
					"Phone": "8(903)755 -49 -29"
				},
				"AccountQuality": 6,
				"Archived": "",
				"ClientId": 15334232,
				"CountryId": 0,
				"CreatedAt": "sa",
				"Currency": "",
				"Grants": [],
				"Login": "flokolibri -	direct",
				"Notification": {
					"ApiObject": {
						"ForceFields": [],
						"NullFields": []
					},
					"NotificationGeneralClients": {
						"ApiObject": {
							"ForceFields": [],
							"NullFields": []
						},
						"Email": "es@123.ru",
						"EmailSubscriptions": [],
						"Lang": "ru"
					}
				},
				"Representatives": [
					{
						"ApiObject": {
							"ForceFields": [],
							"NullFields": []
						},
						"Email": " flokolibri@yandex.ru",
						"Login": "flokolibri -direct",
						"Role": "CHIEF"
					}
				],
				"Restrictions": [],
				"Settings": [],
				"Type": "client",
				"SUBCLIENT": "sad",
				"VatRate": 0
			}]
	}
}`
