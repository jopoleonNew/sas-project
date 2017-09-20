package vkontakte

var Config = struct {
	VKAppID       string
	VKAppSecret   string
	VKRedirectURL string
}{}

func SetParams(AppID, AppSecret, RedirectURL string) {
	Config.VKAppID = AppID
	Config.VKAppSecret = AppSecret
	Config.VKRedirectURL = RedirectURL
}
