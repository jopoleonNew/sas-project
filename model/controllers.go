package model

import (
	"fmt"
	"log"

	"gogs.itcloud.pro/SAS-project/sas/app"

	"gopkg.in/mgo.v2"
)

type Controller struct {
	Config  *app.ConfigType
	Session *mgo.Session
	session *mgo.Session
}

var (
	Config   *app.ConfigType
	DBname   string
	mongoUrl string

	YandexDirectAppID     string
	YandexDirectAppSecret string
)

func Init() {

	//var Config *model.Configuration
	Config = app.GetConfig()
	DBname = Config.DBname
	mongoUrl = Config.Mongourl

	YandexDirectAppID = Config.YandexDirectAppID
	YandexDirectAppSecret = Config.YandexDirectAppSecret
	log.Println("dbcontrollers Init: ", Config)
	return
}

func NewController() (*Controller, error) {
	// This function will return to us a
	// Controller that has our common DB context.
	// We can then use it for multiple routes

	Config = app.GetConfig()
	log.Println("Inside NewController", Config)
	DBname = Config.DBname
	mongoUrl = Config.Mongourl

	YandexDirectAppID = Config.YandexDirectAppID
	YandexDirectAppSecret = Config.YandexDirectAppSecret

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
