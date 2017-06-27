package modelPostgre

import (
	"errors"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

func (u *UserInfo) Create() error {
	// We store only lowercase username
	u.Username = strings.ToLower(u.Username)
	//log.Println("User.AdvanceUpdate() used with ", u)
	if u.Username == "" {
		return errors.New("UserInfo.Update() username field can't be empty")
	}
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(u.collName)
	passbyte := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passbyte, bcrypt.DefaultCost)
	if err != nil {
		log.Println("UserInfo.Create() bcrypt.GenerateFromPassword error: ", err)
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

func (u *UserInfo) AdvanceUpdate() error {

	log.Println("User.AdvanceUpdate() used with ", u)
	if u.Username == "" && u.Email == "" {
		return errors.New("UserInfo.AdvanceUpdate() username and email field can't be empty simultaneously")
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
	if len(u.AccountList) != 0 {
		colQuerier1 := bson.M{"username": u.Username}
		change1 := bson.M{"$push": bson.M{"accountlist": u.AccountList[0]}}
		_, err := c.Upsert(colQuerier1, change1)
		if err != nil {
			log.Println("a.AdvanceUpdate() err: ", err)
			return err
		}
	}
	if len(changeParams) != 0 {
		colQuerier := bson.M{
			"$or": []interface{}{
				bson.M{"email": u.Email},
				bson.M{"username": u.Username},
			},
		}
		change := bson.M{"$set": changeParams}
		log.Println("UserInfo) AdvanceUpdate() query: ", colQuerier, " and cahnge: ", change)
		_, err := c.Upsert(colQuerier, change)
		if err != nil {
			log.Println("a.AdvanceUpdate() err: ", err)
			return err
		}
	}

	return nil
}
