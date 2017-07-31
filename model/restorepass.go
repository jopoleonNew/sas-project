package model

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type RestorePass struct {
	Email     string    `json:"email" bson:"email,omitempty"`
	SecretKey string    `json:"secretkey" bson:"secretkey,omitempty"`
	CreatedAt time.Time `json:"createdat" bson:"createdat,omitempty"`
}

func AddLinkKey(email, linkkey string) error {
	// We store only lowercase username
	//username = strings.ToLower(username)
	//log.Println("User.AdvanceUpdate() used with ", u)

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
		log.Println("AddLinkKey c.Insert(restoreStruct) error: ", err)
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
			log.Println("AddLinkKey c.EnsureIndex(sessionTTL) error: ", err)
			return err
		}
	}
	return nil
}

func MatchKey(linkkey string) (email string, ok bool, err error) {
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C("restoreKeys")
	//err := c.Find()
	//pipeline := bson.M{
	//	"$or": []interface{}{
	//		bson.M{"email": u.Email},
	//		bson.M{"username": u.Username},
	//	},
	//}
	var result RestorePass
	err = c.Find(bson.M{"secretkey": linkkey}).One(&result)
	if err != nil {
		return "", false, err
	}
	if result.SecretKey == linkkey {
		err = c.Remove(bson.M{"secretkey": linkkey})
		if err != nil {
			return "", false, err
		}
		return result.Email, true, nil
	} else {
		return "", false, nil
	}
}
