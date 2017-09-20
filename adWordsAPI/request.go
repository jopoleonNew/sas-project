package adWordsAPI

import (
	"github.com/emiddleton/gads"
)

func Request() {
	authConf := gads.AuthConfig{}
	authConf, err := gads.NewCredentials()
	auth := gads.Auth{
		CustomerId:     "858677780334-qvsl98g1ohh31i0m59oij8ka40p8ne7u.apps.googleusercontent.com",
		DeveloperToken: "FyZsgbxqs_L0RnJbelnh1A",
	}

	campaignService := gads.NewCampaignService(&auth)

	campaigns, totalCount, err := campaignService.Get(
		gads.Selector{
			Fields: []string{
				"Id",
				"Name",
				"Status",
			},
		},
	)
	//
}
