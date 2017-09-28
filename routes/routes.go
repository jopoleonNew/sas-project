package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	//_ "net/http/pprof"

	"gogs.itcloud.pro/SAS-project/sas/handlers/user"
	"gogs.itcloud.pro/SAS-project/sas/handlers/yandex"
	mv "gogs.itcloud.pro/SAS-project/sas/routes/middleware"
)

func LoadRoutes() *mux.Router {

	//making new gorilla router
	r := mux.NewRouter()

	// Bunch of new handlers with gorilla mux router
	//
	// getting auth links for redirecting user in browser to give permision to his account
	// on {source} service
	r.HandleFunc("/getauthlink/{source}", mv.CheckIsUserLogged(user.GetAuthLink))
	// receives requests from {source} services with auth tokens in it
	r.HandleFunc("/gettoken/{source}", mv.CheckIsUserLogged(user.GetToken))

	// collect statistic for one account
	//TODO:  spread logic in two places like this:
	//TODO: mux := mux.NewRouter() mux.Handle("/", myHandler(???)).Methods("GET")
	// returns HTML template to user
	r.HandleFunc("/getaccountstat/{source}", mv.CheckIsUserLogged(user.GetAccountStat)).Methods("GET")
	// receives requests from browser of logged in user to return him statistic of account
	r.HandleFunc("/getaccountstat/{source}", mv.CheckIsUserLogged(user.GetAccountStat)).Methods("POST")
	r.HandleFunc("/isloggedin", user.IsLoggedIn).Methods("GET")
	//r.HandleFunc("/getstatistic/{source}", CheckIsUserLogged(user.GetStatistic))

	// parse template with input forms for {accountlogin} account
	r.HandleFunc("/reports/{accountlogin}",
		mv.CheckIsUserLogged(user.ReportTemplateHandler))

	r.HandleFunc("/", user.IndexHandler) // GET

	// handlers for registration and restoring password
	r.HandleFunc("/signup", user.SignUpHandler)              // GET
	r.HandleFunc("/signupsubmit", user.SignUpSubmitHandler)  // POST
	r.HandleFunc("/activateuser", user.ActivateUserHandler)  // POST
	r.HandleFunc("/forgetpass", user.ForgetPassHandler)      // POST
	r.HandleFunc("/restorepass", user.RestorePasswordHadler) // POST
	//http.HandleFunc("/changepass", user.ChangePasswordHandler)  // POST

	r.HandleFunc("/loginsubmit", user.LoginSubmitHandler)
	r.HandleFunc("/logoutsubmit", mv.CheckIsUserLogged(user.LogoutSubmitHandler))

	r.HandleFunc("/accounts", mv.CheckIsUserLogged(user.AccountsHandler))
	r.HandleFunc("/accountsdata", mv.CheckIsUserLogged(user.AccountsData))
	//http.HandleFunc("/addaccount", user.AddAccountHandler)
	r.HandleFunc("/deleteaccount", user.DeleteAccountHandler)
	// old version
	//r.HandleFunc("/getyandexaccesstoken", yandex.GetYandexAccessToken)

	r.HandleFunc("/getcampaingstats", yandex.GetCampaingStatsHandler) //POST
	r.HandleFunc("/refreshdbcampaign", yandex.RefreshCampaignsListHandler)
	r.HandleFunc("/getreport", yandex.GetStatSliceHandler)
	r.HandleFunc("/fullreport", yandex.ReportTemplateHandler)

	//serving static files for templates
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	return r

}
