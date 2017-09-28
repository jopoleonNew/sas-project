package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/shared/config"
)

func LogRequset(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var store = sessions.NewCookieStore([]byte(config.GetConfig().SessionSecret))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Origin, Content-Type, X-Auth-Token, Authorization, Username, Password")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		//w.Header().Set("Access-Control-Allow-Origin", "https://"+r.Host+"/loginsubmit"+", http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
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
