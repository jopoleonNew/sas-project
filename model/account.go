package model

import (
	//"errors"
	"errors"
	"log"

	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	//Username of user created this account
	Username string
	//Source is the name of organization which hosts the account.
	Source string
	//Accountlogin is the login of account in organization from Source
	Accountlogin       string
	Email              string
	SsaAppYandexID     string
	SsaAppYandexSecret string
	Status             string
	OauthToken         string
	YandexRole         string `json:"yandexrole" bson:"yandexrole"`
	AgencyClients      []string
	CampaignsInfo      []Campaign `json:"campaignsinfo" bson:"campaignsinfo"`
	collName           string     `json:"-"` // mgo Collection name
}

func NewAccount() *Account {
	a := new(Account)
	a.collName = "accountsList"
	return a
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
	// if a.Email == "" {
	// 	return errors.New("Account's Email can't be blank.")
	// }
	// if a.SsaAppYandexID == "" {
	// 	return errors.New("Account's SsaAppYandexID can't be blank.")
	// }
	// if a.SsaAppYandexSecret == "" {
	// 	return errors.New("Account's SsaAppYandexSecret can't be blank.")
	// }
	// if a.Status == "" {
	// 	return errors.New("Account's Status can't be blank.")
	// }
	// if a.OauthToken == "" {
	// 	return errors.New("Account's OauthToken can't be blank.")
	// }
	return nil
}

//// Update updates Account struct fields in
//// database according to passed account as method receiver.
//// Currently DEPRECATED method, use AdvanceUpdate() instead
//func (a *Account) Update() error {
//	log.Println("account.Update used")
//	s := mainSession.Clone()
//	defer s.Close()
//	c := s.DB(mainDB.Name).C(a.collName)
//	err := a.checkMainFields()
//	if err != nil {
//		//log.Println("a.Update checkFields() err: ", err)
//		return err
//	}
//
//	changeInfo, err := c.Upsert(bson.M{"username": a.Username, "accountlogin": a.Accountlogin}, a)
//	if err != nil {
//		log.Println("a.Update Update() err: ", err)
//		return err
//	}
//
//	log.Printf("\n Account ", a, " Upserted in database: %+v ", changeInfo)
//
//	return nil
//}
//
//// SetStatusAndToken() sets Status And Token of account in DB.
//// Currently DEPRECATED method, use AdvanceUpdate() instead
//func (a *Account) SetStatusAndToken() error {
//	err := a.checkMainFields()
//	if err != nil {
//		return err
//	}
//	s := mainSession.Clone()
//	defer s.Close()
//	c := s.DB(mainDB.Name).C(a.collName)
//
//	colQuerier := bson.M{"username": a.Username, "accountlogin": a.Accountlogin, "source": "Яндекс Директ"}
//	change := bson.M{"$set": bson.M{"status": "active", "oauthtoken": a.OauthToken}}
//	err = c.Update(colQuerier, change)
//	if err != nil {
//		log.Println("a.SetStatusAndToken() Update() err: ", err)
//		return err
//	}
//	log.Printf("\n Account ", a, " Upserted in database: %+v ")
//	return nil
//}
//GetInfoList returns slice of Account model with info about accounts from db by username.
// Currently DEPRECATED method, use User.GetAccountList() instead
//func (a *Account) GetInfoList() ([]Account, error) {
//	log.Println("GetInfoList used by ", a.Username, " ", a.Accountlogin)
//	s := mainSession.Clone()
//	defer s.Close()
//	c := s.DB(mainDB.Name).C(a.collName)
//	collen, _ := c.Count()
//	result := make([]Account, collen)
//	err := c.Find(bson.M{"username": a.Username}).All(&result)
//
//	if err != nil {
//		//log.Println(err, "GetInfoList")
//
//		return result, errors.New("Some field in Account are empty")
//	}
//	return result, nil
//}
