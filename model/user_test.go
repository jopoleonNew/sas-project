package model

//
//import (
//	"log"
//	"testing"
//
//	"gopkg.in/mgo.v2/bson"
//)
//
////2017/05/25 14:39:07 UserInfo GetInfo used by  test
////20000           1745700 ns/op
////PASS
////ok      gogs.itcloud.pro/SAS-project/sas/model  52.529s
//// 3700 accounts,  search without indexing on 20 account in AccSlist of User
////user.AccountList = []string{"tbttd", "tbrbe", "taert", "tcbac", "tdggc", "ttcac", "ttdta", "trtga", "tdtge", "tbtgc", "tgccd", "tbdre", "ttbgb", "tgtab", "tactd", "tdrge", "tcebc", "ttcda", "tcrtr", "tbtee"}
//
////Second test after using
////db.testIndexSearch.createIndex( { "accountlogin": 1 } )
////2017/05/25 15:40:27 UserInfo GetInfo used by  test
////100000            468700 ns/op
////PASS
////ok      gogs.itcloud.pro/SAS-project/sas/model  51.511s
////100000            268340 ns/op
////PASS
////ok      gogs.itcloud.pro/SAS-project/sas/model  29.466s
//func TestUserInfo_Create(t *testing.T) {
//	user := NewUser()
//	user.Username = "test"
//	user.Password = "test"
//	user.Create()
//	s := mainSession.Clone()
//	defer s.Close()
//	c := s.DB(mainDB.Name).C("usersList")
//	result := UserInfo{}
//	err := c.Find(bson.M{"username": "test"}).One(&result)
//	if err != nil {
//		t.Error(err)
//	}
//	expected := "test"
//	if result.Username != expected {
//		t.Fatalf("Expected %v, got %v", expected, result.Username)
//	}
//}
//func TestUserInfo_GetInfo(t *testing.T) {
//	user := NewUser()
//	user.Username = "test"
//	//user.Password = "test"
//	result, err := user.GetInfo()
//	if err != nil {
//		t.Error(err)
//	}
//	expected := "test"
//	if result.Username != expected {
//		t.Fatalf("Expected %v, got %v", expected, result.Username)
//	}
//}
//func TestUserInfo_AdvanceUpdate(t *testing.T) {
//	user := NewUser()
//	user.Username = "test"
//	user.Password = "test"
//	user.Email = "testemail"
//	newAccs := []string{"test1", "test2"}
//	user.AccountList = append(user.AccountList, newAccs...)
//	err := user.AdvanceUpdate()
//	if err != nil {
//		t.Error(err)
//	}
//	result, err := user.GetInfo()
//	if err != nil {
//		t.Error(err)
//	}
//	expected := "testemail"
//	if result.Email != expected {
//		t.Fatalf("Expected %v, got %v", expected, result.Email)
//	}
//	log.Println("Inside TEST result.AccountList: ", result.AccountList)
//	expected2 := []string{"test1", "test2"}
//	if result.AccountList[0] != expected2[0] && result.AccountList[1] != expected2[1] {
//		t.Fatalf("Expected %v, got %v", result.AccountList, expected2)
//	}
//}
//func TestUserInfo_IsExist(t *testing.T) {
//	user := NewUser()
//	user.Username = "test12"
//	exists, err := user.IsExist()
//	if err != nil {
//		t.Error(err)
//	}
//	expected := false
//	if exists != expected {
//		t.Fatalf("Expected %v, got %v", expected, exists)
//	}
//
//	user2 := NewUser()
//	user2.Username = "test"
//	exists, err = user2.IsExist()
//	if err != nil {
//		t.Error(err)
//	}
//	expected = true
//	if exists != expected {
//		t.Fatalf("Expected %v, got %v", expected, exists)
//	}
//}
//func TestUserInfo_IsPasswordValid(t *testing.T) {
//	user := NewUser()
//	user.Username = "TEST"
//	//user.Password = "test"
//	valid, err := user.IsPasswordValid("test")
//	if err != nil {
//		t.Error(err)
//	}
//	expected := true
//	if valid != expected {
//		t.Fatalf("Expected %v, got %v", expected, valid)
//	}
//
//	valid, err = user.IsPasswordValid("badtest")
//	if err != nil {
//		t.Error(err)
//	}
//	expected = false
//	if valid != expected {
//		t.Fatalf("Expected %v, got %v", expected, valid)
//	}
//}
//func TestUserInfo_GetAccountList(t *testing.T) {
//
//	acc := NewAccount()
//	acc.Username = "test"
//}
//func TestUserInfo_RemoveAccount(t *testing.T) {
//
//}
//func BenchmarkUserInfo_GetAccountList(b *testing.B) {
//	user := NewUser()
//	user.Username = "test"
//	user.AccountList = []string{"tbttd", "tbrbe", "taert", "tcbac", "tdggc", "ttcac", "ttdta", "trtga", "tdtge", "tbtgc", "tgccd", "tbdre", "ttbgb", "tgtab", "tactd", "tdrge", "tcebc", "ttcda", "tcrtr", "tbtee"}
//	user.AdvanceUpdate()
//	for i := 0; i < b.N; i++ {
//		_, err := user.GetAccountList()
//		if err != nil {
//			b.Error(err)
//		}
//	}
//}
//
////func BenchmarkUserInfo_GetAccountListV2(b *testing.B) {
////	user := NewUser()
////	user.Username = "test"
////	user.AccountList = []string{"tbttd", "tbrbe", "taert", "tcbac", "tdggc", "ttcac", "ttdta", "trtga", "tdtge", "tbtgc", "tgccd", "tbdre", "ttbgb", "tgtab", "tactd", "tdrge", "tcebc", "ttcda", "tcrtr", "tbtee"}
////	user.Update()
////	for i := 0; i < b.N; i++ {
////		_, err := user.GetAccountListV2()
////		if err != nil {
////			b.Error(err)
////		}
////	}
////}
////func BenchmarkUserInfo_GetAccountListV3(b *testing.B) {
////	user := NewUser()
////	user.Username = "test"
////	user.AccountList = []string{"tbttd", "tbrbe", "taert", "tcbac", "tdggc", "ttcac", "ttdta", "trtga", "tdtge", "tbtgc", "tgccd", "tbdre", "ttbgb", "tgtab", "tactd", "tdrge", "tcebc", "ttcda", "tcrtr", "tbtee"}
////	user.Update()
////	for i := 0; i < b.N; i++ {
////		_, err := user.GetAccountListV3()
////		if err != nil {
////			b.Error(err)
////		}
////	}
////}
