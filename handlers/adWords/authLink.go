package adWords

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func GetAdWordsAuthLink(w http.ResponseWriter, r *http.Request) {
	AdWordsURL :=
		"https://accounts.google.com/o/oauth2/auth?" +
			"scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadwords" +
			"&include_granted_scopes=true&prompt=consent&state=example_state_Egor&redirect_uri=" +
			Config.AdWordsRedirectURL +
			"&response_type=code&client_id=" +
			Config.AdWordsAppID + "&access_type=offline"
		///gettoken/adwords?scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadwords&redirect_uri=https%3A%2F%2Fdevelopers.google.com%2Foauthplayground&response_type=code&client_id=858677780334-o9j53qkj1u06p1fafs86gjeu3co8rd0n.apps.googleusercontent.com&access_type=offline

		//https://accounts.google.com/o/oauth2/v2/auth?redirect_uri=https%3A%2F%2Fdevelopers.google.com%2Foauthplayground&prompt=consent&response_type=code&client_id=407408718192.apps.googleusercontent.com&scope=&access_type=offline

		//https://accounts.google.com/o/oauth2/auth?scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadwords&redirect_uri=https%3A%2F%2Fdevelopers.google.com%2Foauthplayground&response_type=code&client_id=858677780334-o9j53qkj1u06p1fafs86gjeu3co8rd0n.apps.googleusercontent.com&access_type=offline

		//https://1cbe3c38.eu.ngrok.io/gettoken/adwords?scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadwords&redirect_uri=https%3A%2F%2Fdevelopers.google.com%2Foauthplayground&response_type=code&client_id=858677780334-o9j53qkj1u06p1fafs86gjeu3co8rd0n.apps.googleusercontent.com&access_type=offline
	logrus.Info("GetAdWordsAuthLink LINK: \n", AdWordsURL)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//adWordsAPI.Request()
	w.Write([]byte(AdWordsURL))
}
