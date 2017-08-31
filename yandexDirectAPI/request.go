package yad

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func (a *Account) Request(url string, fieldNames []string) ([]byte, error) {
	req := &RequestV5{
		Method: "get",
		ParamsV5: &ParamsV5{
			SelectionCriteria: struct{}{},
			FieldNames:        fieldNames,
		},
	}

	reqbytes, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("json.Marshal error: " + err.Error())
	}

	client := &http.Client{}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqbytes))
	if err != nil {
		return nil, errors.New("http.NewRequest error: " + err.Error())
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept-Language", "ru")
	r.Header.Add("Client-Login", a.Login)
	r.Header.Add("Authorization", "Bearer "+a.OAuthToken)
	r.Header.Add("Client-ID", application.ID)

	//log.Println("makeV5GetRequest body Request: ", r)

	resp, err := client.Do(r)

	if err != nil {
		return nil, errors.New("http.Request.Do error: " + err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ioutil.ReadAll error: " + err.Error())
	}
	if string(body) == "" {
		return nil, errors.New("makeV5GetRequest body response is empty")
	}
	//log.Println("makeV5GetRequest Response: ", string(body))
	if strings.Contains(string(body), "error") {
		var errresult = YandexV5Error{}
		err := json.Unmarshal(body, &errresult)
		if err != nil {
			return nil, errors.New("makeV5GetRequest json.Unmarshal error:  " + err.Error())
		}
		return body, errors.New("makeV5GetRequest Yandex Direct API error: " + errresult.ErrorCode + errresult.ErrorString + errresult.ErrorDescription + errresult.RequsetID + string(body))
	}
	//log.Println("makeV5GetRequest body response: ", string(body))
	return body, nil
}

type RequestV4 struct {
	Method  string   `json:"method"`
	Token   string   `json:"token"`
	ParamV4 *ParamV4 `json:"param"`
}

type ParamV4 struct {
	CampaignIDS []int  `json:"CampaignIDS"`
	StartDate   string `json:"StartDate"`
	EndDate     string `json:"EndDate"`
}

type RequestV5 struct {
	Method   string    `json:"method"`
	ParamsV5 *ParamsV5 `json:"params"`
}

type ParamsV5 struct {
	SelectionCriteria struct{} `json:"SelectionCriteria"`
	FieldNames        []string `json:"FieldNames"`
}
