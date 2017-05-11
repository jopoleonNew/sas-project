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
	// TODO: ensure that such notation for pipeline is correct
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

// IsPasswordValid checks is passed password string equals hashed password from DB
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

//
//func (ctl *Controller) AddUserToDB(username, pass, email, name, organ string) error {
//
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//	passbyte := []byte(pass)
//	hashedPassword, err := bcrypt.GenerateFromPassword(passbyte, bcrypt.DefaultCost)
//	if err != nil {
//		log.Println("AddUserToDB bcrypt.GenerateFromPassword error: ", err)
//		return err
//	}
//	usersListCollecton := dbsession.DB(DBname).C("usersList")
//	err = usersListCollecton.Insert(
//		&model.UserInfo{
//			Username:      username,
//			Password:      pass,
//			Salt:          string(hashedPassword),
//			Email:         email,
//			Name:          name,
//			Role:          "regular",
//			Organization:  organ,
//			Registred:     "true",
//			IsActivated:   "false",
//			ActivationKey: "123",
//		})
//	if err != nil {
//		log.Println("AddUserToDB usersListCollecton.Insert error: ", err)
//		return err
//	}
//	log.Println("User ", username, " added to database")
//
//	return nil
//}
//
//func (ctl *Controller) GetUserInfo(username string) (model.UserInfo, error) {
//
//	log.Println("GetUserIfno used")
//	log.Println("GetUserIfno used DBname ", DBname)
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//	c := dbsession.DB(DBname).C("usersList")
//	result := model.UserInfo{}
//	err := c.Find(bson.M{"username": username}).One(&result)
//	if err != nil {
//		log.Println(err, "GetUserIfno")
//		return result, err
//	}
//	return result, nil
//}
//
//func (ctl *Controller) IsUserUnique(username string) (bool, error) {
//	log.Println("IsUserUnique used")
//	userinfo, err := ctl.GetUserInfo(username)
//	if err != nil {
//		log.Println("IsUserUnique GetUserIfno error: ", err)
//		return false, err
//	}
//	if userinfo.Username != "" {
//		return false, nil
//	}
//	return true, nil
//}
//
//func (ctl *Controller) IsUserRegistered(username string) bool {
//	//log.Print("IsRegistred username type: ", reflect.TypeOf(username))
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//	c := dbsession.DB(DBname).C("usersList")
//	result := model.UserInfo{}
//	err := c.Find(bson.M{"username": username}).One(&result)
//	if err != nil {
//		log.Println(err, "IsUserRegistered")
//		return false
//	}
//	if result.Registred == "true" {
//		return true
//	} else {
//		return false
//	}
//}
//
//func (ctl *Controller) IsUserActivated(username string) bool {
//	//log.Print("IsActive username type: ", reflect.TypeOf(username))
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//	c := dbsession.DB(DBname).C("usersList")
//	result := model.UserInfo{}
//	err := c.Find(bson.M{"username": username}).One(&result)
//	if err != nil {
//		log.Println(err, "IsUserActivated")
//		return false
//	}
//	if result.IsActivated == "true" {
//		return true
//	} else {
//		return false
//	}
//}
//
//func (ctl *Controller) ValidateUserPassword(username, password string) (bool, error) {
//	dbsession := ctl.session.Clone()
//	defer dbsession.Close()
//	c := dbsession.DB(DBname).C("usersList")
//	result := model.UserInfo{}
//	err := c.Find(bson.M{"username": username}).One(&result)
//	if err != nil {
//		log.Println("ValidateUserPassword IsUserRegistered error: ", err)
//		return false, err
//	}
//	//hash passwrod+result.salt
//	//хешируем с солью
//	err = bcrypt.CompareHashAndPassword([]byte(result.Salt), []byte(password))
//	if err != nil {
//		return false, err
//	} else {
//		return true, nil
//	}
//	// if result.Password == password {
//	// 	return true, nil
//	// } else {
//	// 	return false, nil
//	// }
//}
//
//func (ctl *Controller) ChangeUsersPassword(username, newpassword string) {
//
//}
