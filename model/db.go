package model

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

var mainSession *mgo.Session
var mainDB mgo.Database
var err error

func SetDBParams(url, dbname string) error {
	info := mgo.DialInfo{
		Addrs:    []string{url},
		Database: dbname,
	}
	log.Println("SetDBParams used")

	mainSession, err = mgo.DialWithInfo(&info)
	if err != nil {

		mainSession.Close()
		log.Fatal("SetDBParams DialWithInfo fatal error: ", err)
		return err
	}
	mainDB = mgo.Database{
		Name:    dbname,
		Session: mainSession,
	}
	if err = mainSession.Ping(); err != nil {
		mainSession.Close()
		log.Fatal("SetDBParams mainSession.Ping() fatal error:", err)
		return err
	}

	log.Println("SetDBParams params: ", mainDB)
	return nil
}
