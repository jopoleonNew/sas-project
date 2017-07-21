package model

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type Account2 struct {
	//Username of user created this account
	Creator string `json:"creator" bson:"creator"`
	//Source is the name of organization which hosts the account.
	Source string `json:"source" bson:"source"`
	//Accountlogin is the login or id of account in organization from Source
	Accountlogin string `json:"accountlogin" bson:"accountlogin"`
	//Owners is the list of user's who have access to that account
	Owners []string `json:"owners" bson:"owners"`

	Email string `json:"email" bson:"email"`
	//active or notactive
	Status string `json:"status" bson:"status"`

	// auth token to make request to source API
	AuthToken string `json:"authtoken" bson:"authtoken"`

	AppID     string `json:"appid" bson:"appid"`
	AppSecret string `json:"appsecret" bson:"appsecret"`

	Role          string     `json:"role" bson:"role"`
	AccountType   string     `json:"accounttype" bson:"accounttype"`
	AgencyClients []string   `json:"agencyclients" bson:"agencyclients"`
	CampaignsInfo []Campaign `json:"campaignsinfo" bson:"campaignsinfo"`
	CreatedAt     time.Time  `json:"createdat" bson:"createdat"`
	LastUpdated   time.Time  `json:"lastupdated" bson:"lastupdated"`
	collName      string     `json:"-"` // mgo Collection name
}

func NewAccount2(Creator, Source, Accountlogin, Email string) *Account2 {
	return &Account2{
		Creator:      Creator,
		Source:       Source,
		Accountlogin: Accountlogin,
		Email:        Email,
		collName:     "accountsList",
	}
}

func (a *Account2) AdvanceUpdate() error {

	//log.Printf("Account.AdvanceUpdate() used with %+v", a)
	err := a.checkMainFields()
	if err != nil {
		return fmt.Errorf("Account.AdvanceUpdate() error: %v", err)
	}

	var changeParams = []bson.DocElem{}
	if a.Creator != "" {
		changeParams = append(changeParams, bson.DocElem{"creator", a.Creator})
	}
	if a.Email != "" {
		changeParams = append(changeParams, bson.DocElem{"email", a.Email})
	}
	if a.AppID != "" {
		changeParams = append(changeParams, bson.DocElem{"appid", a.AppID})
	}
	if a.AppSecret != "" {
		changeParams = append(changeParams, bson.DocElem{"appsecret", a.AppSecret})
	}
	if a.Status != "" {
		changeParams = append(changeParams, bson.DocElem{"status", a.Status})
	}
	if a.AuthToken != "" {
		changeParams = append(changeParams, bson.DocElem{"authtoken", a.AuthToken})
	}
	if a.Role != "" {
		changeParams = append(changeParams, bson.DocElem{"role", a.Role})
	}
	if a.AccountType != "" {
		changeParams = append(changeParams, bson.DocElem{"accounttype", a.AccountType})
	}
	if a.CreatedAt.String() != "0001-01-01 00:00:00 +0000 UTC" {
		changeParams = append(changeParams, bson.DocElem{"createdat", a.CreatedAt})
	}
	if len(a.AgencyClients) != 0 {
		changeParams = append(changeParams, bson.DocElem{"agencyclients", a.AgencyClients})
	}
	if len(a.CampaignsInfo) != 0 {
		changeParams = append(changeParams, bson.DocElem{"campaignsinfo", a.CampaignsInfo})
	}
	changeParams = append(changeParams, bson.DocElem{"lastupdated", time.Now()})
	if len(changeParams) == 0 {
		return errors.New("Nothing to update")
	}
	a.Accountlogin = strings.ToLower(a.Accountlogin)
	a.Creator = strings.ToLower(a.Creator)
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)

	colQuerier := bson.M{"accountlogin": a.Accountlogin, "source": a.Source}

	change := bson.M{"$set": changeParams}
	if len(a.Owners) != 0 {
		colQuerier := bson.M{"accountlogin": a.Accountlogin, "source": a.Source}
		change1 := bson.M{"$push": bson.M{"owners": a.Owners[0]}}
		_, err := c.Upsert(colQuerier, change1)
		if err != nil {
			logrus.Errorf("a.AdvanceUpdate() err: %+v", err)
			return err
		}
	}
	//omitting changeInfo value
	_, err = c.Upsert(colQuerier, change)
	if err != nil {
		logrus.Errorf("a.AdvanceUpdate() err: %+v", err)
		return err
	}
	return nil
}

func (a *Account2) GetInfo() (Account2, error) {

	//log.Println("GetInfo used by ", a.Username, " ", a.Accountlogin)
	if a.Accountlogin == "" {
		return Account2{}, fmt.Errorf("GetInfo() error: Some field in Account are empty")
	}
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)

	result := Account2{}
	err := c.Find(bson.M{"accountlogin": a.Accountlogin}).One(&result)
	if err != nil {
		log.Println("GetAccountIfno err: ", err)
		return Account2{}, err
	}
	return result, nil
}

func (a *Account2) GetAccountList() ([]Account2, error) {

	if a.Accountlogin == "" {
		return []Account2{}, fmt.Errorf("GetAccountList() error: Some field in Account are empty")
	}
	result := []Account2{}
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)
	err = c.Find(bson.M{"accountlogin": bson.M{"$in": a.Owners}}).All(&result)
	if err != nil {
		logrus.Printf("a.GetAccountList() c.Find with login %s, error: %v", a.Accountlogin, err)
		return result, err
	}
	return result, nil
}

//Remove() removes user who called it from Owners field inside account
func (a *Account2) Remove() error {
	log.Println("Remove used: ", a.Creator, " ", a.Accountlogin)
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)
	colQuerier := bson.M{"accountlogin": a.Accountlogin, "source": a.Source}
	change1 := bson.M{"$pull": bson.M{"owners": a.Accountlogin}}
	_, err := c.Upsert(colQuerier, change1)
	if err != nil {
		logrus.Errorf("a.Remove() err: %+v", err)
		return fmt.Errorf("Account.Remove c.Upsert error: %v", err)
	}
	if err != nil {
		logrus.Errorf("Account.Remove error: %v", err)
		return fmt.Errorf("Account.Remove error: %v", err)
	}
	return nil
}

func (a *Account2) checkMainFields() error {
	if a.Source == "" {
		return errors.New("Account's Source field can't be blank.")
	}
	if a.Accountlogin == "" {
		return errors.New("Account's Accountlogin field can't be blank.")
	}

	return nil
}
