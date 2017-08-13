package main

import (
	"log"

	"gogs.itcloud.pro/SAS-project/sas/app"
	"gogs.itcloud.pro/SAS-project/sas/model"
)

func main() {
	app.InitConf("conf-docker.json")

	Config = app.GetConfig()

	log.Printf("CONFIG FILE MAIN: %+v", Config)
	// initiation of MongoDB session
	err := model.SetDBParams(Config.Mongourl, Config.DBname)
	if err != nil {
		log.Fatal(err)
	}
}
func makeNewAccount(index string) error {
	return nil
}
