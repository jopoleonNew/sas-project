package session

import "github.com/gorilla/sessions"

func GetSession(secret string) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(secret))
}
