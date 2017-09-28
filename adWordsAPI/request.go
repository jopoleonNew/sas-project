package adWordsAPI

import (
	gads "github.com/emiddleton/gads"
	"github.com/sirupsen/logrus"
)

//Sas-Web-API
//Идентификатор клиента
//858677780334-o9j53qkj1u06p1fafs86gjeu3co8rd0n.apps.googleusercontent.com
//Секрет клиента
//p8qRT8axYAgwNB145QYK6K0g
//Дата создания
//22 сент. 2017 г., 18:38:14
//client-ID 823426693752-29amuq7jmn1n6eo0qg8ol2tko259i5pt.apps.googleusercontent.com
//client-secret kymFJKVckRLdK9O1i-Ctfmbq
//type Auth struct {
//	CustomerId     string
//	DeveloperToken string
//	UserAgent      string
//	PartialFailure bool
//	Testing        *testing.T   `json:"-"`
//	Client         *http.Client `json:"-"`
//}
// AdWords jopoleon manager client ID :  426-020-5484
// Developer token : FyZsgbxqs_L0RnJbelnh1A
// jopoleon
// Менеджер  • ID: 426-020-5484
//

var Creds = `{
    "oauth2.Config": {
        "ClientID": "546292264965-iie1pr8m7v4l4of303uo9bb17ajdccmu.apps.googleusercontent.com",
        "ClientSecret": "LMHbtEuUhFGhMx8nXvKs1O2f",
        "Endpoint": {
            "AuthURL": "https://accounts.google.com/o/oauth2/auth",
            "TokenURL": "https://accounts.google.com/o/oauth2/token"
        },
        "RedirectURL": "urn:ietf:wg:oauth:2.0:oob",
        "Scopes": [
            "https://adwords.google.com/api/adwords"
        ]
    },
    "oauth2.Token": {
        "access_token": "ya29.GlvPBCtg_KeXdiADnsAOgRM9KNPuuYNU9Sv7SY6AzidjgZ3mCXFKxCjGnPzW0-QOs0KpkQowe4jz5EMkgHo347efCcJ2L7cmfE1OuLpx2E9Zs_6QTUa17S5xF59i",
        "token_type": "Bearer",
        "refresh_token": "1/8xfnonQbaQXrIFR2xpGMrJgaD3rzieCM-R112SUAnn0",
        "expiry": "2017-09-23T22:21:06.1050112+03:00"
    },
    "gads.Auth": {
        "CustomerId": "426-020-5484",
        "DeveloperToken": "FyZsgbxqs_L0RnJbelnh1A",
        "UserAgent": "tests (Golang 1.4 github.com/emiddleton/gads)",
        "PartialFailure": false
    }
}`
var C = gads.Credentials{
	Config: gads.OAuthConfigArgs{
		ClientID:     "546292264965-iie1pr8m7v4l4of303uo9bb17ajdccmu.apps.googleusercontent.com",
		ClientSecret: "LMHbtEuUhFGhMx8nXvKs1O2f",
	},
	Token: gads.OAuthTokenArgs{
		AccessToken:  "ya29.GlvPBCtg_KeXdiADnsAOgRM9KNPuuYNU9Sv7SY6AzidjgZ3mCXFKxCjGnPzW0",
		RefreshToken: "1/8xfnonQbaQXrIFR2xpGMrJgaD3rzieCM-R112SUAnn0",
	},
	Auth: gads.Auth{
		CustomerId:     "426-020-5484",
		DeveloperToken: "FyZsgbxqs_L0RnJbelnh1A",
		UserAgent:      "tests (Golang 1.4 github.com/emiddleton/gads)",
		PartialFailure: false,
	},
}

func Request() {
	authConf, err := gads.NewCredentialsFromParams(C)
	if err != nil {
		logrus.Error("adWordsAPI Request gads.NewCredentialsFromFile error:", err)
		return
	}

	campaignService := gads.NewCampaignService(&authConf.Auth)

	campaigns, totalCount, err := campaignService.Get(
		gads.Selector{
			Fields: []string{
				"Id",
				"Name",
				"Status",
			},
		},
	)

	if err != nil {
		logrus.Error("adWordsAPI Request campaignService.Get error:", err)
		return
	}
	logrus.Infof("Result from Adwords: \n %+v \n %+v", campaigns, totalCount)
	//
}
