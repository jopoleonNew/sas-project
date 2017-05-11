package model

import (
	"log"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type Members struct {
	Name string
	Role string
}

type TemplateInfoStruct struct {
	CurrentUser string
	UsingReport string
	AccountList []Account
	FullStats   []GetSummaryStatRes
}

type Configuration struct {
	Mongourl        string `json:"mongourl"`
	DBname          string `json:"bdname"`
	Yandexappid     string `json:"yandexappid"`
	Yandexappsecret string `json:"yandexappsecret"`
	Yandexapiurl    string `json:"yandexapiurl"`
	Sessionsecret   string `json:"sessionsecret"`
	Serverport      string `json:"serverport"`
	MongoSession    *mgo.Session
}

type GetSummaryStat struct {
	Token  string `json:"token"`
	Method string `json:"method"`
	Param  struct {
		CampaignIDS string `json:"CampaignIDS"`
		StartDate   string `json:"StartDate"`
		EndDate     string `json:"EndDate"`
	} `json:"param"`
}

type Campaign struct {
	ID     int    `json:"Id"`
	Name   string `json:"Name"`
	Status string `json:"Status"`
}

type ModelCampaigns []struct {
	Camp Campaign
}

type CampaignsType struct {
	ID     int    `json:"Id"`
	Name   string `json:"Name"`
	Status string `json:"Status"`
}
type ResultType struct {
	Campaigns []CampaignsType `json:"Campaigns"`
}
type CampaingsGetResult struct {
	Result ResultType `json:"result"`
}

// YandexTokenbody is used in MakeYandexOauthRequest()
// to unmarshal yandex response body and get AccessToken
type YandexTokenbody struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
type SelectionCriteriaType struct{}

type ParamsType struct {
	SelectionCriteria SelectionCriteriaType `json:"SelectionCriteria"`
	FieldNames        []string              `json:"FieldNames"`
}
type GetCampaingsList struct {
	Method string     `json:"method"`
	Params ParamsType `json:"params"`
}
type AgencyInfo struct {
	Result struct {
		Clients []struct {
			Login           string                `json:"login"`
			Representatives []RepresentativesType `json:"representatives"`
		} `json:"clients"`
	} `json:"result"`
}
type RepresentativesType struct {
	Email string `json:"email"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

// "CampaignID": (int),
//        "StatDate": (date),
//        "SumSearch": (float),
//        "SumContext": (float),
//        "ShowsSearch": (int),
//        "ShowsContext": (int),
//        "ClicksSearch": (int),
//        "ClicksContext": (int),
//        "SessionDepthSearch": (float),
//        "SessionDepthContext": (float),
//        "GoalConversionSearch": (float),
//        "GoalConversionContext": (float),
//        "GoalCostSearch": (float),
//        "GoalCostContext": (float)
//  SessionDepthSearch    string  `json:"SessionDepthSearch"`
// SumSearch             int     `json:"SumSearch"`
// ClicksContext         int     `json:"ClicksContext"`
// SessionDepthContext   string  `json:"SessionDepthContext"`
// StatDate              string  `json:"StatDate"`
// GoalCostSearch        float32 `json:"GoalCostSearch"`
// GoalConversionContext string  `json:"GoalConversionContext"`
// ShowsContext          int     `json:"ShowsContext"`
// SumContext            int     `json:"SumContext"`
// GoalConversionSearch  string  `json:"GoalConversionSearch"`
// ShowsSearch           int     `json:"ShowsSearch"`
// CampaignID            int     `json:"CampaignID"`
// GoalCostContext       float32 `json:"GoalCostContext"`
// ClicksSearch          int     `json:"ClicksSearch"`
// interface{}

type YandexTime struct {
	Time time.Time
}

const ctLayout = "2006-01-02"

func (ct *YandexTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	ct.Time, err = time.Parse(ctLayout, string(b))
	if err != nil {
		log.Fatal("models UnmarshalJSON time.Parse error: ", err)
		return err
	}
	return nil
}

func (ct *YandexTime) MarshalJSON() ([]byte, error) {
	// if ct.Time.UnixNano() == nilTime {
	// 	return []byte("null"), nil
	// }
	return []byte(strconv.Quote(ct.Time.Format(ctLayout))), nil
}

type StatisticDataType struct {
	SessionDepthSearch    interface{} `json:"SessionDepthSearch"`
	SumSearch             float32     `json:"SumSearch"`
	ClicksContext         int         `json:"ClicksContext"`
	SessionDepthContext   interface{} `json:"SessionDepthContext"`
	StatDate              YandexTime  `json:"StatDate"`
	GoalCostSearch        interface{} `json:"GoalCostSearch"`
	GoalConversionContext interface{} `json:"GoalConversionContext"`
	ShowsContext          interface{} `json:"ShowsContext"`
	SumContext            interface{} `json:"SumContext"`
	GoalConversionSearch  interface{} `json:"GoalConversionSearch"`
	ShowsSearch           interface{} `json:"ShowsSearch"`
	CampaignID            int         `json:"CampaignID"`
	GoalCostContext       interface{} `json:"GoalCostContext"`
	ClicksSearch          int         `json:"ClicksSearch"`
}
type GetSummaryStatRes struct {
	Data []StatisticDataType `json:"Data"`
}

//implementing Sort.sort interface for GetSummaryStatRes struct
func (p GetSummaryStatRes) Len() int {
	return len(p.Data)
}

func (p GetSummaryStatRes) Less(i, j int) bool {
	return p.Data[i].StatDate.Time.Before(p.Data[j].StatDate.Time)
}

func (p GetSummaryStatRes) Swap(i, j int) {
	p.Data[i], p.Data[j] = p.Data[j], p.Data[i]
}

//////

//////

type GetCampaingsStatisticType struct {
	Token  string    `json:"token"`
	Method string    `json:"method"`
	Param  ParamType `json:"param"`
}
type ParamType struct {
	CampaignIDS []string `json:"CampaignIDS"`
	StartDate   string   `json:"StartDate"`
	EndDate     string   `json:"EndDate"`
}
