//

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"
)

type ConfigType struct {
	Mongourl string `json:"mongourl"`
	DBname   string `json:"bdname"`

	YandexDirectAppID     string `json:"yandexappid"`
	YandexDirectAppSecret string `json:"yandexappsecret"`
	YandexDirectAPIURL    string `json:"yandexapiurl"`

	VKAppID       string `json:"vkappid"`
	VKAppSecret   string `json:"vkappsecret"`
	VKRedirectURL string `json:"vkredirecturl"`

	AdWordsAppID       string `json:"adwordsappid"`
	AdWordsAppSecret   string `json:"adwordsappsecret"`
	AdWordsRedirectURL string `json:"adwordsredirecturl"`

	SessionSecret string       `json:"sessionsecret"`
	ServerPort    string       `json:"serverport"`
	Session       *mgo.Session `json:"-"`
}

var cfg = new(ConfigType)

func GetConfig() *ConfigType {
	return cfg
}

func InitConf(filename string) error {
	//log.Println("....InitConf used")
	var c = &ConfigType{
		Mongourl:              "",
		YandexDirectAppID:     "",
		YandexDirectAppSecret: "",
		YandexDirectAPIURL:    "",
		SessionSecret:         "",
		ServerPort:            "",
	}
	if filename == "" {
		return fmt.Errorf(`Error: Don't dicribe path to a config file.`)
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err = genConfig(filename); err != nil {
			return err
		}
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Can't read config file: %s", err)
	}
	log.Println("cfg file " + filename + " found. Reading...")

	if err = json.Unmarshal(data, c); err != nil {
		logrus.Errorf("Can't read config file: %s", err)
		return err
	}
	cfg = c
	return nil
}

func genConfig(filename string) error {
	logrus.Warn("NOPE. There is no such config file ", filename)
	logrus.Warn("Configuration file not found. Created new with name " + filename + ". " +
		"\n 		     Please, fill it with values you need and RESTART application")
	f, err := os.Create(filename)
	if err != nil {
		log.Println("ReadConfigFile os.Create error: ", err)
		return err
	}
	f.Close()

	var initjson = ConfigType{
		Mongourl:      "localhost",
		DBname:        "sas",
		SessionSecret: "secret",
	}

	writebytes, err := json.MarshalIndent(initjson, "", "\t")
	if err != nil {
		panic(err)
		return err
	}
	err = ioutil.WriteFile(filename, writebytes, 0644)
	if err != nil {
		panic(err)
		return err
	}
	os.Exit(5)
	return nil
}
