package main

import (
	"gogs.itcloud.pro/SAS-project/sas/app"
	userhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/user"
	yandexhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/yandex"
	"gogs.itcloud.pro/SAS-project/sas/model"
	yad "gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
	//"model"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/context"
	//"old_userhandlers"
	_ "net/http/pprof"
)

var (
	//BDctl  *dbcontrollers.Controller
	Config *app.ConfigType
)

func init() {
	var configFileName string

	flag.StringVar(&configFileName, "config", "configuration.json",
		"Specify configuration file name to use. File should be in folder you starting the application")

	flag.Parse()

	app.InitConf(configFileName)

	Config = app.GetConfig()

	log.Printf("CONFIG FILE MAIN: %+v", Config)
	err := model.SetDBParams(Config.Mongourl, Config.DBname)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Debug1")
	err = yad.SetParams(Config.YandexDirectAppID, Config.YandexDirectAppSecret)
	if err != nil {
		log.Fatal(err)
	}
	userhandlers.SetParams(Config)
	yandexhandlers.SetParams(Config)
	//log.Println("Debug2")

}

func main() {

	log.Println("............main() Main.go")
	http.HandleFunc("/", userhandlers.IndexHandler) // GET

	http.HandleFunc("/signup", userhandlers.SignUpHandler)              // GET
	http.HandleFunc("/signupsubmit", userhandlers.SignUpSubmitHandler)  // POST
	http.HandleFunc("/activateuser", userhandlers.ActivateUserHandler)  // POST
	http.HandleFunc("/forgetpass", userhandlers.ForgetPassHandler)      // POST
	http.HandleFunc("/restorepass", userhandlers.RestorePasswordHadler) // POST
	//http.HandleFunc("/changepass", userhandlers.ChangePasswordHandler)  // POST

	http.HandleFunc("/loginsubmit", userhandlers.LoginSubmitHandler)
	http.HandleFunc("/logoutsubmit", userhandlers.LogoutSubmitHandler)

	http.HandleFunc("/accounts", userhandlers.AccountsHandler)
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

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Println("Server started at port: " + Config.ServerPort)
	err := http.ListenAndServe(":"+Config.ServerPort, context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatalln(err)
	}

}
