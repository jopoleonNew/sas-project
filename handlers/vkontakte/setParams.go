package vkontakte

import (
	"gogs.itcloud.pro/SAS-project/sas/app"
)

var Config *app.ConfigType

func SetParams(config *app.ConfigType) {
	Config = config
}
