package yad

var oauthTokenURL = "https://oauth.yandex.ru/token"

const (
	sandBoxURL = "https://api-sandbox.direct.yandex.ru"
	mainAPIURL = "https://api.direct.yandex.ru"
)

var apiURL = mainAPIURL

type Application struct {
	ID     string
	Secret string
}

var application Application

func SetParams(id, secret string)  {
	application.ID = id
	application.Secret = secret

}
