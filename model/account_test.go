package model

import (
	"log"
	"testing"

	"gogs.itcloud.pro/SAS-project/sas/app"
	"gopkg.in/mgo.v2/bson"
)

func init() {

	app.InitConf("../configuration.json")
	var Config = app.GetConfig()
	log.Printf("TESTING CONFIG FILE MAIN: %+v", Config)
	err := SetDBParams(Config.Mongourl, Config.DBname)
	if err != nil {
		log.Fatal(err)
	}

}
func TestNewAccount(t *testing.T) {
	acc := NewAccount()
	if acc.collName != "accountsList" {
		t.Fatalf("Expected %v, got %v", "accountsList", acc.collName)
	}
}

func TestCheckMainFields(t *testing.T) {
	acc := NewAccount()
	acc.Username = "test"
	err := acc.checkMainFields()
	if err == nil {
		t.Error("Expecet err, but there is no")
	}
	acc.Accountlogin = "test"
	err = acc.checkMainFields()
	if err == nil {
		t.Error("Expecet err, but there is no")
	}
	acc.Source = "test"
	err = acc.checkMainFields()
	if err != nil {
		t.Error("There must not be an error, but here it is: ", err)
	}
}

func TestAccount_AdvanceUpdate(t *testing.T) {
	acc := NewAccount()
	acc.Username = "test"
	acc.Accountlogin = "test"
	acc.Source = "test"
	acc.Email = "test"
	acc.YandexRole = "test"
	acc.SsaAppYandexID = "test"
	acc.SsaAppYandexSecret = "test"
	acc.Status = "test"
	acc.OauthToken = "test"
	acc.AgencyClients = []string{"test"}
	testcamp := []Campaign{}
	acc.CampaignsInfo = testcamp
	err = acc.AdvanceUpdate()
	if err != nil {
		t.Error(err)
	}
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C("accountsList")
	result := Account{}
	err := c.Find(bson.M{"accountlogin": "test"}).One(&result)
	if err != nil {
		t.Error(err)
	}
	expected := "test"
	if result.YandexRole != expected {
		t.Fatalf("Expected %v, got %v", expected, result.YandexRole)
	}
	if result.Status != expected {
		t.Fatalf("Expected %v, got %v", expected, result.Status)
	}
}

func TestAccount_Remove(t *testing.T) {
	acc := NewAccount()
	acc.Username = "test"
	acc.Accountlogin = "test"
	acc.Source = "test"
	acc.YandexRole = "test"
	err = acc.Remove()
	if err != nil {
		t.Error(err)
	}
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C("accountsList")
	err := c.Find(bson.M{"accountlogin": "test"}).One(nil)
	if err.Error() != "not found" {

		t.Error(err)

	}
}

func TestAccount_IsExist(t *testing.T) {
	acc := NewAccount()
	acc.Username = "test"
	acc.Accountlogin = "test"
	acc.Source = "test"
	acc.YandexRole = "test"
	//err = acc.Remove()
	//if err != nil {
	//	t.Error(err)
	//}
	exist, err := acc.IsExist()
	if err == ErrAccNotFound {

	} else {
		t.Error(err)
	}
	expected := false
	if exist != expected {
		t.Fatalf("Expected %v, got %v", expected, exist)
	}
	err = acc.AdvanceUpdate()
	if err != nil {
		t.Error(err)
	}
	exist2, err := acc.IsExist()
	log.Println(err)
	if err != nil {
		if err != ErrAccNotFound {
			t.Error(err)
		}
	}
	expected2 := true
	if exist2 != expected2 {
		t.Fatalf("Expected %v, got %v", expected2, exist2)
	}
	err = acc.Remove()
	if err != nil {
		t.Error(err)
	}
}

func TestAccount_GetInfo(t *testing.T) {
	acc := NewAccount()
	acc.Username = "test"
	acc.Accountlogin = "test"
	acc.Source = "test"
	acc.Email = "test"
	acc.YandexRole = "test"
	acc.SsaAppYandexID = "test"
	acc.SsaAppYandexSecret = "test"
	acc.Status = "test"
	acc.OauthToken = "test"
	acc.AgencyClients = []string{"test"}
	testcamp := []Campaign{}
	acc.CampaignsInfo = testcamp
	err = acc.AdvanceUpdate()
	if err != nil {
		t.Error(err)
	}
	result, err := acc.GetInfo()
	if err != nil {
		t.Error(err)
	}
	expected := "test"
	if result.YandexRole != expected {
		t.Fatalf("Expected %v, got %v", expected, result.YandexRole)
	}
	if result.Status != expected {
		t.Fatalf("Expected %v, got %v", expected, result.Status)
	}
}

func BenchmarkAccount_GetInfo(b *testing.B) {
	acc := NewTestAccount()
	acc.Username = "test"
	acc.Accountlogin = "tcrtr"
	acc.Source = "test"
	for i := 0; i < b.N; i++ {
		_, err := acc.GetInfo()
		if err != nil {
			b.Error(err)
		}
	}

}

//func BenchmarkAccount_IsExist40(b *testing.B) {
//	acc := NewAccount()
//	acc.Username = "test"
//	acc.Accountlogin = "test"
//	acc.Source = "test"
//	acc.YandexRole = "test"
//	err = acc.AdvanceUpdate()
//	if err != nil {
//		b.Error(err)
//	}
//	for i := 0; i < b.N; i++ {
//		_, err := acc.IsExist()
//		if err != nil {
//			b.Error(err)
//		}
//	}
//	err = acc.Remove()
//	if err != nil {
//		b.Error(err)
//	}
//
//}
