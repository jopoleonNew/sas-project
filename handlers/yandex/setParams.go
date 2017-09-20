package yandex

//var Config *configuration.ConfigType
var Config = struct {
	YandexDirectAppID     string
	YandexDirectAppSecret string
}{}

func SetParams(AppID, AppSecret string) {

	Config.YandexDirectAppID = AppID
	Config.YandexDirectAppSecret = AppSecret

}
