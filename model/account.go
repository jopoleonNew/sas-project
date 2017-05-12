package model

import (
	//"errors"
	"errors"
	"log"

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

var ErrAccNotFound = errors.New("Account is not exist.")

// IsExist checks DB for existing account with given account login.
func (a *Account) IsExist() (bool, error) {
	log.Println("IsExist used")
	//if a.Username == "" {
	//	return false, errors.New("IsExist Account's Username field can't be blank.")
	//}
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

// Update updates Account struct fields in
// database according to passed account as method receiver.
// TODO: 1. Check duplication fields of inserting json's. (???)
// TODO: 2. Check for already deleted Yandex campaigns but yet saved in DB. (DONE_???)
func (a *Account) Update() error {
	log.Println("account.Update used")
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)
	err := a.checkMainFields()
	if err != nil {
		//log.Println("a.Update checkFields() err: ", err)
		return err
	}

	changeInfo, err := c.Upsert(bson.M{"username": a.Username, "accountlogin": a.Accountlogin}, a)
	if err != nil {
		log.Println("a.Update Update() err: ", err)
		return err
	}

	log.Printf("\n Account ", a, " Upserted in database: %+v ", changeInfo)

	return nil
}

func (a *Account) SetStatusAndToken() error {
	err := a.checkMainFields()
	if err != nil {
		return err
	}
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)

	colQuerier := bson.M{"username": a.Username, "accountlogin": a.Accountlogin, "source": "Яндекс Директ"}
	change := bson.M{"$set": bson.M{"status": "active", "oauthtoken": a.OauthToken}}
	err = c.Update(colQuerier, change)
	if err != nil {
		log.Println("a.SetStatusAndToken() Update() err: ", err)
		return err
	}
	log.Printf("\n Account ", a, " Upserted in database: %+v ")
	return nil
}

func (a *Account) AdvanceUpdate() error {

	log.Println("Account.AdvanceUpdate() used with ", a)
	err := a.checkMainFields()
	if err != nil {
		return err
	}

	var changeParams = []bson.DocElem{}

	//if a.Username != "" {
	//	changeParams = append(changeParams, bson.DocElem{"username", a.Username})
	//}
	//if a.Accountlogin != "" {
	//	changeParams = append(changeParams, bson.DocElem{"accountlogin", a.Accountlogin})
	//}
	//if a.Source != "" {
	//	changeParams = append(changeParams, bson.DocElem{"source", a.Source})
	//}
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
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)

	colQuerier := bson.M{"username": a.Username, "accountlogin": a.Accountlogin, "source": a.Source}
	change := bson.M{"$set": changeParams}
	//log.Println("......................................AdvanceUpdate queries :    ", colQuerier, changeParams)
	_, err = c.Upsert(colQuerier, change)
	if err != nil {
		log.Println("a.AdvanceUpdate() err: ", err)
		return err
	}
	//changeInfostr := fmt.Sprintf("%+v", changeInfo1)
	//log.Printf("\n Account %+v Updated in database ", a)
	return nil
}

func (a *Account) SetCampaigns() error {
	return nil
}

//GetInfo returns Account model with info about account from db by username and account login.
func (a *Account) GetInfo() (Account, error) {

	log.Println("GetInfo used by ", a.Username, " ", a.Accountlogin)
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)
	// err := a.checkMainFields()
	// if err != nil {
	// 	log.Println("GetInfo checkMainFields err: ", err)
	// 	return nil, err
	// }
	if a.Username == "" || a.Accountlogin == "" {
		return Account{}, errors.New("Some field in Account are empty")
	}
	result := Account{}
	err := c.Find(bson.M{"username": a.Username, "accountlogin": a.Accountlogin}).One(&result)
	if err != nil {
		log.Println("GetAccountIfno err: ", err)
		return Account{}, err
	}
	return result, nil
}

//GetInfoList returns slice of Account model with info about accounts from db by username.
func (a *Account) GetInfoList() ([]Account, error) {
	log.Println("GetInfoList used by ", a.Username, " ", a.Accountlogin)
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(a.collName)
	collen, _ := c.Count()
	result := make([]Account, collen)
	err := c.Find(bson.M{"username": a.Username}).All(&result)
	if err != nil {
		//log.Println(err, "GetInfoList")

		return result, errors.New("Some field in Account are empty")
	}
	return result, nil
}

//Remove removes given Account from DB
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

/////////////////////////////****************//////////////////////
/////////////////////////////****************//////////////////////
/////////////////////////////****************//////////////////////
/////////////////////////////****************//////////////////////

//GetAccountInfo returns info about account with given username and account login.
//func (ctl *Controller) GetAccountInfo(username, accountlogin string) (model.Account, error) {
//
//	log.Println("GetAccountIfno used. Recieved info: ", username, accountlogin)
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//	c := dbsession.DB(DBname).C("accountsList")
//	result := model.Account{}
//	err := c.Find(bson.M{"username": username, "accountlogin": accountlogin}).One(&result)
//	if err != nil {
//		log.Println("GetAccountIfno err: ", err)
//		return result, err
//	}
//	return result, nil
//}
//
//func (ctl *Controller) AddAccountToDB(username, source, accountlogin, email, accrole string) error {
//	log.Println("AddAccountToDB used")
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//
//	accountsListCollecton := dbsession.DB(DBname).C("accountsList")
//	err := accountsListCollecton.Insert(
//		&Account{
//			Username:           username,
//			Source:             source,
//			Accountlogin:       accountlogin,
//			Email:              email,
//			SsaAppYandexID:     YandexDirectAppID,
//			SsaAppYandexSecret: YandexDirectAppSecret,
//			YandexRole:         accrole,
//			Status:             "notactive",
//			OauthToken:         "",
//		})
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	log.Println("Account of user ", username, " with organization ", source, " added to database.")
//
//	return nil
//
//}
//
//func (ctl *Controller) GetAcctountsList(username string) ([]model.Account, error) {
//	log.Println("GetAcctountsList used")
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//	c := dbsession.DB(DBname).C("accountsList")
//	len, _ := c.Count()
//	result := make([]model.Account, len)
//	err := c.Find(bson.M{"username": username}).All(&result)
//	if err != nil {
//		log.Println(err, "GetAcctountsList")
//		return result, err
//	}
//	return result, nil
//}
//
//// IsAccountUnique checks uniqueness of combination username + accountlogin
//// Is there account with given username and account? returns true if no, false if yes
//
//// SetAccountTokenandStatusActive setts Status and YANDEX Accsess Token
//// by account login and username
//func (ctl *Controller) SetAccountTokenandStatusActive(username, accountlogin, token string) error {
//	log.Println("SetAccountTokenandStatusActive used : ", username, token)
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//	c := dbsession.DB(DBname).C("accountsList")
//
//	colQuerier := bson.M{"username": username, "accountlogin": accountlogin, "source": "Яндекс Директ"}
//	change := bson.M{"$set": bson.M{"status": "active", "oauthtoken": token}}
//	err := c.Update(colQuerier, change)
//	if err != nil {
//		log.Println("SetAccountStatus err: ", err)
//		return err
//	}
//	return nil
//}
//
//func (ctl *Controller) DeleteAccountfromDB(username, accountlogin string) error {
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//	accountsListCollecton := dbsession.DB(DBname).C("accountsList")
//	err := accountsListCollecton.Remove(bson.M{"username": username, "accountlogin": accountlogin})
//	if err != nil {
//		log.Println("DeleteAccountfromDB", err)
//		return err
//	}
//	return nil
//}
//
//func (ctl *Controller) AddCampaingsToAccount(username, accountlogin string, campingsstruct []*yad.Campaign) error {
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//
//	var yadcs []model.YadCampaign
//	for _, v := range campingsstruct {
//		yadcs = append(yadcs, model.YadCampaign{
//			ID:     v.ID,
//			Name:   v.Name,
//			Status: v.Status,
//		})
//	}
//
//	accountsListCollecton := dbsession.DB(DBname).C("accountsList")
//	colQuerier := bson.M{"username": username, "accountlogin": accountlogin, "source": "Яндекс Директ"}
//	change := bson.M{"$set": bson.M{"campaignsinfo": &yadcs}}
//	err := accountsListCollecton.Update(colQuerier, change)
//	if err != nil {
//		log.Println("AddCampaingsToAccount err: ", err)
//		return err
//	}
//	return nil
//
//}
//
//func (ctl *Controller) AddOrganizationToDB(login, pass, email, name, organ string) error {
//	return nil
//}
