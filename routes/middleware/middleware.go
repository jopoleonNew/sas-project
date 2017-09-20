package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/shared/config"
)

func CheckIsUserLogged(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var store = sessions.NewCookieStore([]byte(config.GetConfig().SessionSecret))
		session, err := store.Get(r, "sessionSSA")
		if err != nil {
			logrus.Error("CheckIsUserLogged store.Get err: ", err.Error())
			fmt.Fprintf(w, "CheckIsUserLogged store.Get err: "+err.Error())
			return
		}
		//log.Println("CheckIsUserLogged middleware values: ", session.Values)

		if session.Values["loggedin"] != nil && session.Values["loggedin"].(string) == "true" &&
			len(session.Values) != 0 {
			//Add data to context
			ctx := context.WithValue(r.Context(), "username", session.Values["username"])
			next.ServeHTTP(w, r.WithContext(ctx))
			//next.ServeHTTP(w, r)
		} else {
			fmt.Fprintf(w, "You are not logged in.")
			//http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		//log.Println("Executing CheckIsUserLogged again")
	})
}
