package adWords

var Config = struct {
	AdWordsAppID       string
	AdWordsAppSecret   string
	AdWordsRedirectURL string
}{}

func SetParams(AppID, AppSecret, RedirectURL string) {
	Config.AdWordsAppID = AppID
	Config.AdWordsAppSecret = AppSecret
	Config.AdWordsRedirectURL = RedirectURL
}
