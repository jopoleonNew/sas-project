package main

import (
	"gogs.itcloud.pro/SAS-project/sas/app"
	"gogs.itcloud.pro/SAS-project/sas/shared/config"
)

var (
	Config *config.ConfigType
)

func main() {
	app.Run()
}
