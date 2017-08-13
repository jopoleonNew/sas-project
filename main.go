package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"

	"gogs.itcloud.pro/SAS-project/sas/app"
	userhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/user"
	vkhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/vkontakte"
	yandexhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/yandex"
	"gogs.itcloud.pro/SAS-project/sas/model"
	yad "gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"

	_ "net/http/pprof"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	//BDctl  *dbcontrollers.Controller
	Config *app.ConfigType
)

func init() {
	var configFileName string
	var local bool
	flag.StringVar(&configFileName, "config", "conf-docker.json",
		"Specify configuration file name to use. File should be in folder you starting the application")
	flag.BoolVar(&local, "local", false,
		"Specify where the app  started")
	flag.Parse()

	app.InitConf(configFileName)

	Config = app.GetConfig()

	log.Printf("CONFIG FILE MAIN: %+v", Config)
	// initiation of MongoDB session
	err := model.SetDBParams(Config.Mongourl, Config.DBname)
	if err != nil {
		panic(err)
	}

	err = yad.SetParams(Config.YandexDirectAppID, Config.YandexDirectAppSecret)
	if err != nil {
		log.Fatal(err)
	}

	userhandlers.SetParams(Config)
	yandexhandlers.SetParams(Config)
	vkhandlers.SetParams(Config)
}

func CheckIsUserLogged(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var store = sessions.NewCookieStore([]byte(app.GetConfig().SessionSecret))
		session, err := store.Get(r, "sessionSSA")
		if err != nil {
			log.Println(errors.New("CheckIsUserLogged store.Get err: " + err.Error()))
			fmt.Fprintf(w, "CheckIsUserLogged store.Get err: "+err.Error())
			return
		}
		log.Println("CheckIsUserLogged middleware values: ", session.Values)

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
		log.Println("Executing CheckIsUserLogged again")
	})
}

func main() {

	//making new gorilla router
	r := mux.NewRouter()

	// Block of new handlers with gorilla mux router
	//
	// getting auth links for redirecting user in browser to give permision to his account
	// on {source} service
	r.HandleFunc("/getauthlink/{source}", CheckIsUserLogged(userhandlers.GetAuthLink))
	//
	r.HandleFunc("/gettoken/{source}", CheckIsUserLogged(userhandlers.GetToken))

	// collect statistic for one account
	//TODO:  spread logic in two places like this:
	//TODO: mux := mux.NewRouter() mux.Handle("/", myHandler(???)).Methods("GET")
	r.HandleFunc("/getaccountstat/{source}", CheckIsUserLogged(userhandlers.GetAccountStat))
	//r.HandleFunc("/getstatistic/{source}", CheckIsUserLogged(userhandlers.GetStatistic))

	// parse template with input forms for {accountlogin} account
	r.HandleFunc("/reports/{accountlogin}", CheckIsUserLogged(userhandlers.ReportTemplateHandler))

	r.HandleFunc("/", userhandlers.IndexHandler) // GET

	r.HandleFunc("/signup", userhandlers.SignUpHandler)              // GET
	r.HandleFunc("/signupsubmit", userhandlers.SignUpSubmitHandler)  // POST
	r.HandleFunc("/activateuser", userhandlers.ActivateUserHandler)  // POST
	r.HandleFunc("/forgetpass", userhandlers.ForgetPassHandler)      // POST
	r.HandleFunc("/restorepass", userhandlers.RestorePasswordHadler) // POST
	//http.HandleFunc("/changepass", userhandlers.ChangePasswordHandler)  // POST

	r.HandleFunc("/loginsubmit", userhandlers.LoginSubmitHandler)
	r.HandleFunc("/logoutsubmit", userhandlers.LogoutSubmitHandler)

	r.HandleFunc("/accounts", CheckIsUserLogged(userhandlers.AccountsHandler2))
	//http.HandleFunc("/addaccount", userhandlers.AddAccountHandler)
	r.HandleFunc("/deleteaccount", userhandlers.DeleteAccountHandler)

	//TODO: join this two endpoints in one, to make possible autoadding of accounts.
	r.HandleFunc("/getauthcodeyandex", yandexhandlers.GetAuthCodeYandexHandler)
	r.HandleFunc("/submityandexcode", yandexhandlers.SubmitConfirmationYandexCode)

	r.HandleFunc("/getyandexaccesstoken", yandexhandlers.GetYandexAccessToken)

	r.HandleFunc("/getcampaingstats", yandexhandlers.GetCampaingStatsHandler) //POST
	r.HandleFunc("/refreshdbcampaign", yandexhandlers.RefreshCampaignsListHandler)
	r.HandleFunc("/getreport", yandexhandlers.GetStatSliceHandler)
	r.HandleFunc("/fullreport", yandexhandlers.ReportTemplateHandler)

	//"/getauthcodevk", vkhandlers.GetVKAuthCode uses for getting VKapp info from server
	r.HandleFunc("/getauthcodevk", vkhandlers.GetVKAuthCode)
	r.HandleFunc("/vkauth", vkhandlers.VKauthorize)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Println("Server started at port: " + Config.ServerPort)
	//srv := &http.Server{
	//	Handler:      r,
	//	Addr:         "127.0.0.1:8000",
	//	// Good practice: enforce timeouts for servers you create!
	//	WriteTimeout: 15 * time.Second,
	//	ReadTimeout:  15 * time.Second,
	//}
	//
	//log.Fatal(srv.ListenAndServe())
	//err := http.ListenAndServe(":"+Config.ServerPort, gctx.ClearHandler(http.DefaultServeMux))
	err := http.ListenAndServe(":"+Config.ServerPort, r)
	if err != nil {
		log.Fatalln(err)
	}

}
