package utils

import (
	"errors"
	"log"
	"net/http"

	"gogs.itcloud.pro/SAS-project/sas/app"

	"github.com/gorilla/sessions"
)

//var store = sessions.NewCookieStore([]byte("secret"))

var ErrorEmptySessionUsername = errors.New("GetUsernamefromRequestSession error: session.Values['username'] is empty")

// IsUserLoggedIn checking gorilla session user is logged or not
func IsUserLoggedIn(r *http.Request) bool {
	var store = sessions.NewCookieStore([]byte(app.GetConfig().SessionSecret))
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println(err)

	}

	if session.Values["loggedin"] != nil && session.Values["loggedin"].(string) == "true" {
		return true
	} else {
		return false
	}

}

// getUsernamefromRequestSession returns username from incoming request gorilla session
func GetUsernamefromRequestSession(r *http.Request) (string, error) {
	var store = sessions.NewCookieStore([]byte(app.GetConfig().SessionSecret))
	//ErrorEmptySessionUsername.
	//var ErrorString =
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		log.Println(err)
		return "", errors.New("GetUsernamefromRequestSession store.Get err: " + err.Error())
	}

	if session.Values["username"] != nil && session.Values["username"].(string) != "" {
		return session.Values["username"].(string), nil
	}

	return "", ErrorEmptySessionUsername

}
