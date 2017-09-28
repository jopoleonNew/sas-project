package database

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

var Session *mgo.Session
var mainDB mgo.Database
var err error

func SetDBParams(url, dbname string) error {
	info := mgo.DialInfo{
		Addrs:    []string{url},
		Database: dbname,
		Timeout:  3 * time.Second,
	}
	//log.Println("SetDBParams of MongoDB used")

	Session, err = mgo.DialWithInfo(&info)
	if err != nil {
		log.Fatal("SetDBParams DialWithInfo fatal error: ", err)
		Session.Close()
		return err
	}
	//mainSession.
	//err =
	//if err != nil {
	//	mainSession.Close()
	//	log.Fatal("SetDBParams mainSession.Ping() fatal error:", err)
	//	return err
	//}
	mainDB = mgo.Database{
		Name:    dbname,
		Session: Session,
	}
	//log.Println("SetDBParams params: ", mainDB)
	return nil
}
