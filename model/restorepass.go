package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type RestorePass struct {
	Email     string    `json:"email" bson:"email,omitempty"`
	SecretKey string    `json:"secretkey" bson:"secretkey,omitempty"`
	CreatedAt time.Time `json:"createdat" bson:"createdat,omitempty"`
}

// AddLinkKey creates new activation key document to restoreKeys collection with 5 hours TTL
func AddLinkKey(email, linkkey string) error {
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C("restoreKeys")
	restoreStruct := &RestorePass{
		Email:     email,
		SecretKey: linkkey,
		CreatedAt: time.Now(),
	}
	err := c.Insert(restoreStruct)
	if err != nil {
		logrus.Error("AddLinkKey c.Insert(restoreStruct) error: ", err)
		return err
	}
	sessionTTL := mgo.Index{
		Key:         []string{"createdat"},
		Unique:      false,
		DropDups:    false,
		Background:  true,
		ExpireAfter: (5 * time.Hour)} // session_expire is a time.Duration
	//
	err = c.EnsureIndex(sessionTTL)
	if err != nil {
		if err.Error() != "Index with name: createdat_1 already exists with different options" {
			return nil
		} else {
			logrus.Error("AddLinkKey c.EnsureIndex(sessionTTL) error: ", err)
			return err
		}
	}
	return nil
}

// MatchKey matches given key by user's email with created activation key
func MatchKey(linkkey string) (email string, ok bool, err error) {
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C("restoreKeys")
	var result RestorePass
	err = c.Find(bson.M{"secretkey": linkkey}).One(&result)
	if err != nil {
		return "", false, err
	}
	if result.SecretKey == linkkey {
		// removing matching key, if it matches with given
		err = c.Remove(bson.M{"secretkey": linkkey})
		if err != nil {
			return "", false, err
		}
		return result.Email, true, nil
	} else {
		return "", false, nil
	}
}
