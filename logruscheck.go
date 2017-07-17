package main

import (
	api "github.com/nk2ge5k/goyad"
	"github.com/nk2ge5k/yandex-direct-api-go/campaigns"
	// "os"
	"fmt"
	"log"
)

//type GetResponse struct {
//	support.ApiObject
//	general.GetResponseGeneral
//
//	Campaigns *[]CampaignGetItem `json:"Campaigns,omitempty"`
//}

//type CampaignGetItem struct {
//	support.ApiObject
//	CampaignBase
//
//	Id                  *int64                      `json:"Id,omitempty"`
//	Name                *string                     `json:"Name,omitempty"`
//	StartDate           *string                     `json:"StartDate,omitempty"`
//	Type                *CampaignTypeGetEnum        `json:"Type,omitempty"`
//	Status              *general.StatusEnum         `json:"Status,omitempty"`
//	State               *CampaignStateGetEnum       `json:"State,omitempty"`
//	StatusPayment       *CampaignStatusPaymentEnum  `json:"StatusPayment,omitempty"`
//	StatusClarification *string                     `json:"StatusClarification,omitempty"`
//	SourceId            *int64                      `json:"SourceId,omitempty"`
//	Statistics          *general.Statistics         `json:"Statistics,omitempty"`
//	Currency            *general.CurrencyEnum       `json:"Currency,omitempty"`
//	Funds               *FundsParam                 `json:"Funds,omitempty"`
//	RepresentedBy       *CampaignAssistant          `json:"RepresentedBy,omitempty"`
//	DailyBudget         *DailyBudget                `json:"DailyBudget,omitempty"`
//	EndDate             *string                     `json:"EndDate,omitempty"`
//	NegativeKeywords    *general.ArrayOfString      `json:"NegativeKeywords,omitempty"`
//	BlockedIps          *general.ArrayOfString      `json:"BlockedIps,omitempty"`
//	ExcludedSites       *general.ArrayOfString      `json:"ExcludedSites,omitempty"`
//	TextCampaign        *TextCampaignGetItem        `json:"TextCampaign,omitempty"`
//	MobileAppCampaign   *MobileAppCampaignGetItem   `json:"MobileAppCampaign,omitempty"`
//	DynamicTextCampaign *DynamicTextCampaignGetItem `json:"DynamicTextCampaign,omitempty"`
//	TimeTargeting       *TimeTargeting              `json:"TimeTargeting,omitempty"`
//}
//
func main() {
	//"access_token":
	//api.Client{Login: "f20.ru", Token: ""}
	//client := api.NewClient("f20.ru", "API_TOKEN")
	client := api.NewClient()
	//clients.New()"f20.ru", "API_TOKEN"
	client.Token = api.Token{Value: "AQAAAAATYgojAARqcat-1xGYAkGcizp9tXSfi7E"}
	client.Login = "Ctdirect"
	request := campaigns.GetRequest{
		FieldNames: []campaigns.CampaignFieldEnum{

			"ClientInfo",
			"Currency",
			"DailyBudget",
			"EndDate",
			"Id",
			"Name",
			"Notification",
			"RepresentedBy",
			"SourceId",
			"StartDate",
			"State",
			"Statistics",
			"Status",
			"TimeZone",
			"Type",
		},
		TextCampaignFieldNames: []campaigns.TextCampaignFieldEnum{
			"BiddingStrategy",
			"CounterIds",
			"RelevantKeywords",
			"Settings",
		},
	}

	service := campaigns.New(&client)

	result, err := service.Get(request)
	if err != nil {
		log.Fatal(err)
	}
	s := result.Campaigns
	for _, c := range s {

		fmt.Printf("1. %+v\n", c.Name)
		//fmt.Printf("2. %v\n", &c.Name)
	}
	//p := &s

	//fmt.Printf("2. %+v\n", s)
	//fmt.Printf("3. %#v\n", s)
	//fmt.Printf("4. %p\n", &s)
	//fmt.Printf("4. %q\n", &s)
	//fmt.Printf("5. %v\n", &s)
	//logrus.Printf("%v %p", *result.Campaigns)
	//logrus.Printf("Result from Yandex about %s \n %+v", client.Login, result)
}
