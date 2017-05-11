package model

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

var mainSession *mgo.Session
var mainDB mgo.Database

func SetDBParams(url, dbname string) (err error) {
	info := mgo.DialInfo{
		Addrs:    []string{url},
		Database: dbname,
	}

	mainSession, err = mgo.DialWithInfo(&info)
	if err != nil {
		//println(err)
		return err
	}
	mainDB = mgo.Database{
		Name:    dbname,
		Session: mainSession,
	}
	log.Println("SetDBParams params: ", mainDB)
	return nil
}
