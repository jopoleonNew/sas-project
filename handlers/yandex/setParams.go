package yandex

import (
	"gogs.itcloud.pro/SAS-project/sas/app"

	"github.com/gorilla/sessions"
)

var Config *app.ConfigType
var store *sessions.CookieStore

func SetParams(config *app.ConfigType) {

	Config = config

	store = sessions.NewCookieStore([]byte(config.SessionSecret))
}
