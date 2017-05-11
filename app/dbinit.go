package app

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

// func init() {

// 	//var Config *models.Configuration
// 	Config = GetConfig()
// 	// DBname = Config.DBname
// 	// mongoUrl = Config.Mongourl

// 	// YandexDirectAppID = Config.YandexDirectAppID
// 	// YandexDirectAppSecret = Config.YandexDirectAppSecret
// 	log.Println("app/dbinit.go init: ", Config)
// 	return
// }

type Controller struct {
	//Config  *ConfigType
	Session *mgo.Session
	session *mgo.Session
}

// func MakeDBSession() {

// }
// func GetMongoSession() *mgo.Session {

// }

func NewController() (*Controller, error) {
	// This function will return to us a
	// Controller that has our common DB context.
	// We can then use it for multiple routes
	//log.Println("dbcontroller NewController uri to bd: ", Config.Mongourl)
	uri := Config.Mongourl
	if uri == "" {
		return nil, fmt.Errorf("NewController error: no DB connection string provided")
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("NewController mgo.Dial error: " + err.Error())
	}
	return &Controller{
		Session: session,
		session: session,
	}, nil
}
