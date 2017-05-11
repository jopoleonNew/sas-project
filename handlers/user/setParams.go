package user

import (
	"github.com/gorilla/sessions"
	"gogs.itcloud.pro/SAS-project/sas/app"
)

var Config *app.ConfigType
var store *sessions.CookieStore

func SetParams(config *app.ConfigType) {

	Config = config

	store = sessions.NewCookieStore([]byte(config.SessionSecret))
}
