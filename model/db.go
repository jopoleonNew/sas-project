package model

import (
	"log"

	"time"

	mgo "gopkg.in/mgo.v2"
)

var mainSession *mgo.Session
var mainDB mgo.Database
var err error

func SetDBParams(url, dbname string) error {
	info := mgo.DialInfo{
		Addrs:    []string{url},
		Database: dbname,
		Timeout:  3 * time.Second,
	}
	log.Println("SetDBParams used")

	mainSession, err = mgo.DialWithInfo(&info)
	if err != nil {
		log.Fatal("SetDBParams DialWithInfo fatal error: ", err)
		mainSession.Close()
		return err
	}
	log.Println("mainSession.Ping(): ", mainSession.Ping())
	//mainSession.
	//err =
	//if err != nil {
	//	mainSession.Close()
	//	log.Fatal("SetDBParams mainSession.Ping() fatal error:", err)
	//	return err
	//}
	mainDB = mgo.Database{
		Name:    dbname,
		Session: mainSession,
	}

	log.Println("SetDBParams params: ", mainDB)
	return nil
}
