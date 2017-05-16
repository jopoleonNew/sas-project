package model

import (
	"errors"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserInfo struct {
	Id            bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Username      string        `json:"username" bson:"username"`
	Password      string        `json:"password" bson:"password"`
	Salt          string        `json:"salt" bson:"salt"`
	Email         string        `json:"email" bson:"email"`
	AccountList   []string      `json:"accountlist" bson:"accountlist"`
	Name          string        `json:"name" bson:"name"`
	Role          string        `json:"role" bson:"role"`
	Organization  string        `json:"organisation" bson:"organisation"`
	Registred     string        `json:"registered" bson:"registered"`
	IsActivated   string        `json:"isActivated" bson:"isActivated"`
	ActivationKey string        `json:"activationKey" bson:"activationKey"`
	collName      string        `json:"-"`
}

// NewUser returns empty *UserInfo.
func NewUser() *UserInfo {
	u := new(UserInfo)
	u.collName = "usersList"
	return u
}

func (u *UserInfo) Update() error {
	// We store only lowercase username
	u.Username = strings.ToLower(u.Username)

	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(u.collName)
	passbyte := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passbyte, bcrypt.DefaultCost)
	if err != nil {
		log.Println("db user.go Update bcrypt.GenerateFromPassword error: ", err)
		return err
	}
	u.Username = strings.ToLower(u.Username)
	u.Salt = string(hashedPassword)
	//u.collName = "usersList"
	changeInfo, err := c.Upsert(bson.M{"username": u.Username}, u)
	if err != nil {
		log.Println("AddUserToDB usersListCollecton.Insert error: ", err)
		return err
	}
	log.Println("User ", u.Username, " Upserted in database: ", changeInfo)

	return nil
}

//IsExist checks is user with given username or email exist in DB already.
func (u *UserInfo) IsExist() (bool, error) {

	u.Username = strings.ToLower(u.Username)

	log.Println("IsExists user used for: ", u.Username)

	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(u.collName)
	u.Username = strings.ToLower(u.Username)
	log.Println("IsExists user used for: ", u.Username)
	pipeline := bson.M{
		"$or": []interface{}{
			bson.M{"email": u.Email},
			bson.M{"username": u.Username},
		},
	}
	err := c.Find(pipeline).One(nil)
	if err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		} else {
			log.Println("a.IsExist err: ", err)
			return false, err
		}
	}

	return true, nil
}

func (u *UserInfo) GetInfo() (UserInfo, error) {
	log.Println("UserInfo GetInfo used by ", u.Username)
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(u.collName)
	u.Username = strings.ToLower(u.Username)
	if u.Username == "" {
		return UserInfo{}, errors.New("Some field in UserInfo are empty")
	}
	result := UserInfo{}
	u.Username = strings.ToLower(u.Username)
	err := c.Find(bson.M{"username": u.Username}).One(&result)
	if err != nil {
		log.Println("GetInfo err: ", err)
		return UserInfo{}, err
	}
	return result, nil
}

func (u *UserInfo) GetAccountList() ([]Account, error) {
	if u.Username == "" {
		return nil, errors.New("UserInfo.GetAccountList() username field can't be blank.")
	}
	log.Println("UserInfo.GetAccountList() by ", u.Username)

	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C("accountsList")
	//u.Username = strings.ToLower(u.Username)
	result := make([]Account, len(u.AccountList))
	//db.accountsList.find({ "accountlogin": { $in: [ "123qwe","123sssa" ] }})
	log.Println("GetAccountList acclist: ", u.AccountList)
	//var useracclist []string
	userinfo, err := u.GetInfo()
	if err != nil {
		log.Println("GetAccountList u.GetInfo() err: ", err)
		return nil, err
	}
	err = c.Find(bson.M{"accountlogin": bson.M{"$in": userinfo.AccountList}}).All(&result)
	if err != nil {
		log.Println("GetAccountList err: ", err)
		return nil, err
	}
	//result := []Account{}
	//u.AccountList
	return result, nil
}

// IsPasswordValid checks if password string equals hashed password from DB
func (u *UserInfo) IsPasswordValid(password string) (bool, error) {
	u.Username = strings.ToLower(u.Username)
	log.Println("ValidateUserPassword GetInfo used by ", u.Username)
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(u.collName)
	result := UserInfo{}
	u.Username = strings.ToLower(u.Username)
	err := c.Find(bson.M{"username": u.Username}).One(&result)
	if err != nil {
		log.Println("IsPasswordValid c.Find error: ", err)
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Salt), []byte(password))

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (u *UserInfo) AdvanceUpdate() error {

	log.Println("User.AdvanceUpdate() used with ", u)
	if u.Username == "" {
		return errors.New("UserInfo.AdvanceUpdate() username field can't be empty")
	}
	var changeParams = []bson.DocElem{}

	if u.Password != "" {
		changeParams = append(changeParams, bson.DocElem{"password", u.Password})
	}
	if u.Salt != "" {
		changeParams = append(changeParams, bson.DocElem{"salt", u.Salt})
	}
	if u.Email != "" {
		changeParams = append(changeParams, bson.DocElem{"email", u.Email})
	}
	if u.Name != "" {
		changeParams = append(changeParams, bson.DocElem{"role", u.Role})
	}
	if u.Organization != "" {
		changeParams = append(changeParams, bson.DocElem{"registred", u.Registred})
	}
	if u.IsActivated != "" {
		changeParams = append(changeParams, bson.DocElem{"isActivated", u.IsActivated})
	}
	if u.ActivationKey != "" {
		changeParams = append(changeParams, bson.DocElem{"activationKey", u.ActivationKey})
	}
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(u.collName)
	if len(changeParams) != 0 {
		colQuerier := bson.M{"username": u.Username}
		change := bson.M{"$set": changeParams}
		log.Println("UserInfo) AdvanceUpdate() query: ", colQuerier, " and cahnge: ", change)
		_, err := c.Upsert(colQuerier, change)
		if err != nil {
			log.Println("a.AdvanceUpdate() err: ", err)
			return err
		}
	}
	if len(u.AccountList) != 0 {
		colQuerier1 := bson.M{"username": u.Username}
		change1 := bson.M{"$push": bson.M{"accountlist": u.AccountList[0]}}
		_, err := c.Upsert(colQuerier1, change1)
		if err != nil {
			log.Println("a.AdvanceUpdate() err: ", err)
			return err
		}
	}
	return nil
}

// RemoveAccount removes account login from AccountList filed of UserInfo struct
func (u *UserInfo) RemoveAccount() error {
	log.Println("User.RemoveAccount() used with ", u)
	if u.Username == "" {
		return errors.New("UserInfo.PushAccount() username field can't be empty")
	}
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(u.collName)
	if len(u.AccountList) != 0 {
		colQuerier1 := bson.M{"username": u.Username}
		// using $pull MongoDB operator to remove account
		change1 := bson.M{"$pull": bson.M{"accountlist": u.AccountList[0]}}
		_, err := c.Upsert(colQuerier1, change1)
		if err != nil {
			log.Println("u.PushAccount() err: ", err)
			return err
		}
	}
	return nil
}
