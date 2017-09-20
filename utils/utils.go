package utils

import (
	"errors"
	"net/http"

	"gogs.itcloud.pro/SAS-project/sas/shared/config"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

var ErrorEmptySessionUsername = errors.New("GetUsernamefromRequestSession error: session.Values['username'] is empty")

// GetUsernamefromRequestSession returns username from incoming request with gorilla session
func GetUsernamefromRequestSession(r *http.Request) (string, error) {
	var store = sessions.NewCookieStore([]byte(config.GetConfig().SessionSecret))
	session, err := store.Get(r, "sessionSSA")
	if err != nil {
		logrus.Error("GetUsernamefromRequestSession store.Get(r, sessionSSA) err", err)
		return "", errors.New("GetUsernamefromRequestSession store.Get err: " + err.Error())
	}

	if session.Values["username"] != nil && session.Values["username"].(string) != "" {
		return session.Values["username"].(string), nil
	}
	return "", ErrorEmptySessionUsername

}
