package main

import (
	"log"
	"math/rand"

	"gogs.itcloud.pro/SAS-project/sas/app"
	"gogs.itcloud.pro/SAS-project/sas/model"
)

const letterBytes = "abcdegrt"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func main() {
	app.InitConf("configuration.json")

	Config := app.GetConfig()

	log.Printf("CONFIG FILE MAIN: %+v", Config)
	err := model.SetDBParams(Config.Mongourl, Config.DBname)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10000; i++ {
		acc := model.NewTestAccount()

		acc.Accountlogin = "t" + RandStringBytes(4)
		acc.Username = "test"
		acc.Source = "test"
		acc.YandexRole = "test23"
		err := acc.AdvanceUpdate()
		if err != nil {
			log.Fatal(err)
		}
	}

}
