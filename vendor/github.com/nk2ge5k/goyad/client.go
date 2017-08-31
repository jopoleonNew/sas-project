package goyad

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/nk2ge5k/goyad/support"
	"github.com/sirupsen/logrus"
)

const (
	API_URL         string = "https://api.direct.yandex.com/json/v5/"
	API_SANDBOX_URL string = "https://api-sandbox.direct.yandex.com/json/v5/"
)

type TokenInterface interface {
	GetToken() string
}

type Token struct {
	Value string
}

func (t Token) GetToken() string {
	return t.Value
}

type ApiRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

//type YandexV5Error struct {
//	ErrorDescription string `json:"error_description"`
//	ErrorCode        string `json:"error_code"`
//	ErrorString      string `json:"error_string"`
//	RequsetID        string `json:"request_id"`
//}
type ApiError struct {
	RequestId   string `json:"request_id"`
	ErrorCode   string `json:"error_code"`
	ErrorString string `json:"error_string"`
	ErrorDetail string `json:"error_detail"`
}

type ClientInterface interface {
	// Do new API request
	Do(s string, m string, p interface{}) ([]byte, error)
}

type Client struct {
	Login             string
	Token             TokenInterface
	ApiUrl            string
	Language          string
	UseOperatorPoints bool
	HttpClient        http.Client
	LastRequestId     string
	LastCallCost      int
	UnitsLeft         int
	UnitsLimit        int
}

func (e ApiError) Error() string {
	return fmt.Sprintf(
		"Error (Request ID %v, Code %v): %v %v",
		e.RequestId,
		e.ErrorCode,
		e.ErrorString,
		e.ErrorDetail,
	)
}

type FakeClientInterface interface {
}

func (c *Client) Do(s string, m string, p support.MappedObjectInterface) ([]byte, error) {

	n, f := p.Get()

	areq := ApiRequest{
		Method: m,
		Params: support.
			New(p).
			Map(n, f),
	}
	b, err := json.Marshal(areq)
	if err != nil {
		return []byte{}, err
	}
	r, err := http.NewRequest(http.MethodPost, c.ApiUrl+s, bytes.NewReader(b))
	if err != nil {
		return []byte{}, err
	}
	r.Header = http.Header{
		"Content-Type":       []string{"application/json", "charset=utf-8"},
		"Authorization":      []string{"Bearer " + c.Token.GetToken()},
		"Client-Login":       []string{c.Login},
		"Accept-Language":    []string{c.Language},
		"Use-Operator-Units": []string{strconv.FormatBool(c.UseOperatorPoints)},
	}
	resp, err := c.HttpClient.Do(r)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	//ioutil.WriteFile("campstestresponse", []byte(*resp), 0666)
	//logrus.Infof("\n\n\n\n\n request full: \n %s ", resp)
	c.LastRequestId = resp.Header["Requestid"][0]

	if nil != resp.Header["Units"] {
		units := strings.Split(resp.Header["Units"][0], "/")

		if v, err := strconv.Atoi(units[0]); err != nil {
			return []byte{}, err
		} else {
			c.LastCallCost = v
		}

		if v, err := strconv.Atoi(units[1]); err != nil {
			return []byte{}, err
		} else {
			c.UnitsLeft = v
		}

		if v, err := strconv.Atoi(units[2]); err != nil {
			return []byte{}, err
		} else {
			c.UnitsLimit = v
		}
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	//logrus.Info(string(res))
	if bytes.Compare(res[:9], []byte("{\"error\":")) == 0 {
		ex := map[string]ApiError{}
		if err := json.Unmarshal(res, &ex); err != nil {
			logrus.Errorf("ERROR OF API_ERROR STRUCT inside goyad value: %s err:%+v", string(res), err)
			return res, err
		} else {
			return res, ex["error"]
		}
	} else {
		return res, nil
	}
}

// May be create without login and token
func NewClient() Client {
	return Client{
		ApiUrl:   API_URL,
		Language: "ru",
	}
}
