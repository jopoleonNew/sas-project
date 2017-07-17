package model

import (
	//"errors"
	"errors"
	"log"

	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SourceInfo struct {
	Accountlogin  string
	AppID         string
	AppSecret     string
	AuthToken     string
	CampaingsInfo []Campaign
	AgencyClients []string
	AccountRole   string
}
type YandexInfo struct {
	AppID         string
	AppSecret     string
	OauthToken    string
	YandexRole    string `json:"yandexrole" bson:"yandexrole"`
	AgencyClients []string
	CampaignsInfo []Campaign `json:"campaignsinfo" bson:"campaignsinfo"`
}

//account_id integer	идентификатор рекламного кабинета.
//account_type//string	тип рекламного кабинета. Возможные значения:
//general — обычный;
//agency — агентский.
//account_status integer, [0,1]	статус рекламного кабинета. Возможные значения:
//1 — активен;
//0 — неактивен.
//access_role string	права пользователя в рекламном кабинете. Возможные значения:
//admin — главный администратор;
//manager — администратор;
//reports — наблюдатель.
type VKInfo struct {
	VKtoken       string
	AccessRole    string
	AccountType   string
	AccountStatus string
}
type Account struct {
	//Username of user created this account
	Username string
	//Source is the name of organization which hosts the account.
	Source string
	//Accountlogin is the login of account in organization from Source
	Accountlogin string
	//Owners is the list of user's who have access to that account
	Owners []string

	Email string

	Status string

	collName           string `json:"-"` // mgo Collection name
	SourceInfo         SourceInfo
	SsaAppYandexID     string
	SsaAppYandexSecret string
	OauthToken         string
	VKtoken            string
	YandexRole         string `json:"yandexrole" bson:"yandexrole"`
	AgencyClients      []string
	CampaignsInfo      []Campaign `json:"campaignsinfo" bson:"campaignsinfo"`
}
type Account2 struct {
	//Username of user created this account
	Username string
	//Source is the name of organization which hosts the account.
	Source string
	//Accountlogin is the login or id of account in organization from Source
	Accountlogin string
	//Owners is the list of user's who have access to that account
	Owners []string `json:"owners" bson:"owners"`

	Email string

	Status string //active or notactive

	AppID         string
	AppSecret     string
	AuthToken     string
	Role          string
	AccountType   string
	AgencyClients []string
	CampaignsInfo []Campaign `json:"campaignsinfo" bson:"campaignsinfo"`

	collName string `json:"-"` // mgo Collection name
}

func (a *Account) Adapter() {

}

type AccountDB interface {
}

func NewAccount() *Account {
	a := new(Account)
	a.collName = "accountsList"
	return a
}
func NewAccount2(Username, Source, Accountlogin, Email string) *Account2 {
	//a := new(Account)
	//a.collName = "accountsList"
	return &Account2{
		Username:     Username,
		Source:       Source,
		Accountlogin: Accountlogin,
		Email:        Email,
		collName:     "accountsList",
	}
}
func NewTestAccount() *Account {
	a := new(Account)
	a.collName = "testIndexSearch"
	return a
}

var ErrAccNotFound = errors.New("Account is not exist.")

// IsExist checks user's AccountList for existing account with given account login.
func (a *Account) IsExist() (bool, error) {
	log.Println("(a *Account) IsExist used")

	if a.Accountlogin == "" {
		return false, errors.New("IsExist Account's Accountlogin field can't be blank.")
	}

	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)

	err := c.Find(bson.M{"accountlogin": a.Accountlogin}).One(nil)
	if err != nil {
		if err == mgo.ErrNotFound {
			return false, ErrAccNotFound
		} else {
			log.Println("Account.IsExist err: ", err)
			return false, err
		}
	}

	return true, nil
}

// AdvanceUpdate() updates account in DB. Upgraded versiof of Update() method.
// Checks all inbound fields of method receiver and updates document in DB appropriate with
// given values.
func (a *Account) AdvanceUpdate() error {

	//log.Printf("Account.AdvanceUpdate() used with %+v", a)
	err := a.checkMainFields()
	if err != nil {
		return err
	}

	var changeParams = []bson.DocElem{}

	if a.Email != "" {
		changeParams = append(changeParams, bson.DocElem{"email", a.Email})
	}
	if a.SsaAppYandexID != "" {
		changeParams = append(changeParams, bson.DocElem{"ssaappyandexid", a.SsaAppYandexID})
	}
	if a.SsaAppYandexSecret != "" {
		changeParams = append(changeParams, bson.DocElem{"ssaappyandexsecret", a.SsaAppYandexSecret})
	}
	if a.Status != "" {
		changeParams = append(changeParams, bson.DocElem{"status", a.Status})
	}
	if a.OauthToken != "" {
		changeParams = append(changeParams, bson.DocElem{"oauthtoken", a.OauthToken})
	}
	if a.VKtoken != "" {
		changeParams = append(changeParams, bson.DocElem{"vktoken", a.VKtoken})
	}
	if a.YandexRole != "" {
		changeParams = append(changeParams, bson.DocElem{"yandexrole", a.YandexRole})
	}
	if len(a.AgencyClients) != 0 {
		changeParams = append(changeParams, bson.DocElem{"agencyclients", a.AgencyClients})
	}
	if len(a.CampaignsInfo) != 0 {
		changeParams = append(changeParams, bson.DocElem{"campaignsinfo", a.CampaignsInfo})
	}

	if len(changeParams) == 0 {
		return errors.New("Account.AdvanceUpdate() error: Nothing to update")
	}
	a.Accountlogin = strings.ToLower(a.Accountlogin)
	a.Username = strings.ToLower(a.Username)
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)

	colQuerier := bson.M{"username": a.Username, "accountlogin": a.Accountlogin, "source": a.Source}

	change := bson.M{"$set": changeParams}
	if len(a.Owners) != 0 {
		colQuerier1 := bson.M{"accountlogin": a.Accountlogin}
		change1 := bson.M{"$push": bson.M{"owners": a.Owners[0]}}
		_, err := c.Upsert(colQuerier1, change1)
		if err != nil {
			log.Println("a.AdvanceUpdate() err: ", err)
			return err
		}
	}
	//omitting changeInfo value
	_, err = c.Upsert(colQuerier, change)
	if err != nil {
		log.Println("a.AdvanceUpdate() err: ", err)
		return err
	}

	//changeInfostr := fmt.Sprintf("%+v", changeInfo1)
	//log.Printf("\n Account %+v Updated in database ", a)
	return nil
}
func (a *Account2) AdvanceUpdate() error {

	//log.Printf("Account.AdvanceUpdate() used with %+v", a)
	err := a.checkMainFields()
	if err != nil {
		return err
	}

	var changeParams = []bson.DocElem{}

	if a.Email != "" {
		changeParams = append(changeParams, bson.DocElem{"email", a.Email})
	}
	if a.AppID != "" {
		changeParams = append(changeParams, bson.DocElem{"ssaappyandexid", a.AppID})
	}
	if a.AppSecret != "" {
		changeParams = append(changeParams, bson.DocElem{"ssaappyandexsecret", a.AppSecret})
	}
	if a.Status != "" {
		changeParams = append(changeParams, bson.DocElem{"status", a.Status})
	}
	if a.AuthToken != "" {
		changeParams = append(changeParams, bson.DocElem{"oauthtoken", a.AuthToken})
	}
	if a.Role != "" {
		changeParams = append(changeParams, bson.DocElem{"yandexrole", a.Role})
	}
	if len(a.AgencyClients) != 0 {
		changeParams = append(changeParams, bson.DocElem{"agencyclients", a.AgencyClients})
	}
	if len(a.CampaignsInfo) != 0 {
		changeParams = append(changeParams, bson.DocElem{"campaignsinfo", a.CampaignsInfo})
	}

	if len(changeParams) == 0 {
		return errors.New("Account.AdvanceUpdate() error: Nothing to update")
	}
	a.Accountlogin = strings.ToLower(a.Accountlogin)
	a.Username = strings.ToLower(a.Username)
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)

	colQuerier := bson.M{"username": a.Username, "accountlogin": a.Accountlogin, "source": a.Source}

	change := bson.M{"$set": changeParams}
	if len(a.Owners) != 0 {
		colQuerier1 := bson.M{"accountlogin": a.Accountlogin}
		change1 := bson.M{"$push": bson.M{"owners": a.Owners[0]}}
		_, err := c.Upsert(colQuerier1, change1)
		if err != nil {
			log.Println("a.AdvanceUpdate() err: ", err)
			return err
		}
	}
	//omitting changeInfo value
	_, err = c.Upsert(colQuerier, change)
	if err != nil {
		log.Println("a.AdvanceUpdate() err: ", err)
		return err
	}

	//changeInfostr := fmt.Sprintf("%+v", changeInfo1)
	//log.Printf("\n Account %+v Updated in database ", a)
	return nil
}

//GetInfo returns Account model with info about account from db by username and account login.
func (a *Account) GetInfo() (Account, error) {

	log.Println("GetInfo used by ", a.Username, " ", a.Accountlogin)
	if a.Username == "" || a.Accountlogin == "" {
		return Account{}, errors.New("GetInfo() error: Some field in Account are empty")
	}
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)

	result := Account{}
	err := c.Find(bson.M{"username": a.Username, "accountlogin": a.Accountlogin}).One(&result)
	if err != nil {
		log.Println("GetAccountIfno err: ", err)
		return Account{}, err
	}
	return result, nil
}

//Remove() removes given Account from DB
func (a *Account) Remove() error {
	log.Println("Remove used: ", a.Username, " ", a.Accountlogin)
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)
	err := c.Remove(bson.M{"username": a.Username, "accountlogin": a.Accountlogin})
	if err != nil {
		log.Println("model Account Remove error: ", err)
		return err
	}
	return nil
}

func (a *Account) checkMainFields() error {
	if a.Username == "" {
		return errors.New("Account's Username field can't be blank.")
	}
	if a.Source == "" {
		return errors.New("Account's Source field can't be blank.")
	}
	if a.Accountlogin == "" {
		return errors.New("Account's Accountlogin field can't be blank.")
	}

	return nil
}
