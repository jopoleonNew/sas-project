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

	gctx "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gogs.itcloud.pro/SAS-project/sas/modelPostgre"
)

var (
	//BDctl  *dbcontrollers.Controller
	Config *app.ConfigType
)

func init() {
	var configFileName string

	flag.StringVar(&configFileName, "config", "conf-docker.json",
		"Specify configuration file name to use. File should be in folder you starting the application")

	flag.Parse()

	app.InitConf(configFileName)

	Config = app.GetConfig()

	log.Printf("CONFIG FILE MAIN: %+v", Config)
	// initiation of MongoDB session
	err := model.SetDBParams(Config.Mongourl, Config.DBname)
	if err != nil {
		panic(err)
	}

	//host = "localhost"
	//port = 5432
	//user = "postgres"
	//password = "qwe"
	//dbname = "test"

	//initiation of PostgreSQL driver
	modelPostgre.SetDBParams("localhost", 5432, "postgres", "qwe", "test")
	if err != nil {
		panic(err)
	}
	//log.Println("Debug1")
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

			fmt.Fprintf(w, "You are not logger in.")
			//http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		log.Println("Executing CheckIsUserLogged again")
	})
}
func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["source"])
}
func main() {
	//finalHandler := http.HandlerFunc(final)
	//
	//http.Handle("/", middlewareOne(middlewareTwo(finalHandler)))
	//http.ListenAndServe(":3000", nil)
	//dexHandler := http.HandlerFunc(userhandlers.IndexHandler)
	r := mux.NewRouter()
	r.HandleFunc("/addaccount/{source}", ArticlesCategoryHandler)

	http.HandleFunc("/", userhandlers.IndexHandler) // GET

	http.HandleFunc("/signup", userhandlers.SignUpHandler)              // GET
	http.HandleFunc("/signupsubmit", userhandlers.SignUpSubmitHandler)  // POST
	http.HandleFunc("/activateuser", userhandlers.ActivateUserHandler)  // POST
	http.HandleFunc("/forgetpass", userhandlers.ForgetPassHandler)      // POST
	http.HandleFunc("/restorepass", userhandlers.RestorePasswordHadler) // POST
	//http.HandleFunc("/changepass", userhandlers.ChangePasswordHandler)  // POST

	http.HandleFunc("/loginsubmit", userhandlers.LoginSubmitHandler)
	http.HandleFunc("/logoutsubmit", userhandlers.LogoutSubmitHandler)

	http.HandleFunc("/accounts", CheckIsUserLogged(userhandlers.AccountsHandler))
	http.HandleFunc("/addaccount", userhandlers.AddAccountHandler)
	http.HandleFunc("/deleteaccount", userhandlers.DeleteAccountHandler)

	//TODO: join this two endpoints in one, to make possible autoadding of accounts.
	http.HandleFunc("/getauthcodeyandex", yandexhandlers.GetAuthCodeYandexHandler)
	http.HandleFunc("/submityandexcode", yandexhandlers.SubmitConfirmationYandexCode)

	http.HandleFunc("/getyandexaccesstoken", yandexhandlers.GetYandexAccessToken)

	http.HandleFunc("/getcampaingstats", yandexhandlers.GetCampaingStatsHandler) //POST
	http.HandleFunc("/refreshdbcampaign", yandexhandlers.RefreshCampaignsListHandler)
	http.HandleFunc("/getreport", yandexhandlers.GetStatSliceHandler)
	http.HandleFunc("/report", yandexhandlers.ReportTemplateHandler)

	//"/getauthcodevk", vkhandlers.GetVKAuthCode uses for getting VKapp info from server
	http.HandleFunc("/getauthcodevk", vkhandlers.GetVKAuthCode)
	http.HandleFunc("/vkauth", vkhandlers.VKauthorize)

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Println("Server started at port: " + Config.ServerPort)
	err := http.ListenAndServe(":"+Config.ServerPort, gctx.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatalln(err)
	}

}
