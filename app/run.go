package app

import (
	"github.com/sirupsen/logrus"
	userhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/user"
	vkhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/vkontakte"
	yandexhandlers "gogs.itcloud.pro/SAS-project/sas/handlers/yandex"
	"gogs.itcloud.pro/SAS-project/sas/model"

	"flag"
	"net/http"

	"gogs.itcloud.pro/SAS-project/sas/handlers/adWords"
	"gogs.itcloud.pro/SAS-project/sas/routes"
	"gogs.itcloud.pro/SAS-project/sas/shared/config"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

// Run starts application, parsing start flags, sets parameters in packages, inits workers pools

func Run() {

	var configFileName string
	var local bool
	flag.StringVar(&configFileName, "config", "conf-docker.json",
		"Specify configuration file name to use. File should be in folder you starting the application")
	flag.BoolVar(&local, "local", false,
		"Specify where the app  started")
	flag.Parse()
	// reading config file in configFileName path
	err := config.InitConf(configFileName)
	if err != nil {
		logrus.Fatal("Run() config.InitConf error: ", err)
	}
	cfg := config.GetConfig()

	logrus.Printf("CONFIG FILE MAIN: %+v", cfg)
	// setting parameters from config file for packages
	SetParams(cfg)

	//initializing API's work pools for queueing requests
	//yad.InitRequestQueue()
	yad.InitPool(5)
	defer yad.YPool.Close()

	//updaters.Start(time.Second * 15)
	logrus.Info("Server started at port: " + cfg.ServerPort)
	// loading http routes
	r := routes.LoadRoutes()

	// starting http listener
	logrus.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))

}

// SetParams sets parameters in all packages
func SetParams(c *config.ConfigType) {

	userhandlers.SetParams(c.SessionSecret)
	yandexhandlers.SetParams(c.YandexDirectAppID, c.YandexDirectAppSecret)
	vkhandlers.SetParams(c.VKAppID, c.VKAppSecret, c.VKRedirectURL)
	adWords.SetParams(c.AdWordsAppID, c.AdWordsAppSecret, c.AdWordsRedirectURL)
	yad.SetParams(c.YandexDirectAppID, c.YandexDirectAppSecret)
	err := model.SetDBParams(c.Mongourl, c.DBname)
	if err != nil {
		logrus.Fatal(err)
	}

}
