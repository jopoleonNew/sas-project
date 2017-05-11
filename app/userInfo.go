package app

import "gopkg.in/mgo.v2/bson"

type UserInfo struct {
	Id            bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Username      string        `json:username bson:username`
	Password      string        `json:password bson:password`
	Salt          string        `json:salt bson:salt`
	Email         string        `json:email bson:email`
	Name          string        `json:name bson:name`
	Role          string        `json:role bson:role`
	Organization  string        `json:organisation bson:organisation`
	Registred     string        `json:registered bson:registered`
	IsActivated   string        `json:isActivated bson:isActivated`
	ActivationKey string        `json:ActivationKey bson:ActivationKey`
}
