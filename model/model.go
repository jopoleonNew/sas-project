package model

import (
	"log"
	"strconv"
	"time"
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

type Campaign struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Owner  string `json:"owner"`
}

// YandexTokenbody is used in MakeYandexOauthRequest()
// to unmarshal yandex response body and get AccessToken

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
