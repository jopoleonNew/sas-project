package user

import (
	"github.com/gorilla/sessions"
)

//var Config *config.ConfigType
var store *sessions.CookieStore

//var store = session.GetSession()

func SetParams(sessionSecret string) {
	//Config = c
	store = sessions.NewCookieStore([]byte(sessionSecret))
}
