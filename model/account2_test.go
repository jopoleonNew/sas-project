package model

import (
	"log"
	"testing"

	"reflect"

	"time"

	"fmt"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/shared/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	testDB  = "testDB"
	cleanDB = make(chan bool)
	done    = make(chan bool)
)

func init() {
	err := config.InitConf("../configuration.json")
	if err != nil {
		logrus.Fatalf("Reading test config file error: %v ", err)
	}
	var Config = config.GetConfig()
	log.Printf("TESTING CONFIG FILE MAIN: %+v", Config)
	//StartMongoDB(*cfg)
	//time.Sleep(5 * time.Second)
	err = SetDBParams(Config.Mongourl, testDB)
	if err != nil {
		log.Fatalf("Test init SetDBParams error: %v", err)
	}

	go func() {
		for {
			<-cleanDB
			fmt.Println("Test DB ", testDB, " droped.")
			cleanUpTest()
			<-done
			close(cleanDB)
			break
		}

	}()
}

func TestNewAccount2(t *testing.T) {
	newAcc := NewAccount2("t", "t", "t", "t")

	if newAcc.Email != "t" {
		t.Fatalf("Wrong email, got %s, expecte %s", newAcc.Email, "t")
	}

	if newAcc.Creator != "t" {
		t.Fatalf("Wrong Creator, got %s, expecte %s", newAcc.Creator, "t")
	}

	if newAcc.Source != "t" {
		t.Fatalf("Wrong Source, got %s, expecte %s", newAcc.Source, "t")
	}

	if newAcc.Accountlogin != "t" {
		t.Fatalf("Wrong Accountlogin, got %s, expecte %s", newAcc.Accountlogin, "t")
	}
	if newAcc.Owners[0] != "t" {
		t.Fatalf("Wrong Owners, got %s, expecte %s", newAcc.Owners[0], "t")
	}
	if newAcc.collName != "accountsList" {
		t.Fatalf("Wrong Accountlogin, got %s, expecte %s", newAcc.Accountlogin, "t")
	}
}
func TestAccount2_Update(t *testing.T) {
	acc := NewAccount2("test", "test", "test", "test")
	acc.Role = "test"
	acc.AgencyClients = []string{"test"}
	testcamp := []Campaign{{ID: 1, Name: "test1"}, {ID: 2, Name: "test2"}}
	acc.CampaignsInfo = testcamp
	err = acc.Update()
	if err != nil {
		t.Fatal(err)
	}
	s := mainSession.Clone()
	defer s.Close()
	c := s.DB(mainDB.Name).C(acc.collName)
	result := Account2{}
	err := c.Find(bson.M{"accountlogin": "test"}).One(&result)
	if err != nil {
		t.Fatal(err)
	}
	//test updating Role field
	acc.Role = "test2"
	err = acc.Update()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Find(bson.M{"accountlogin": "test"}).One(&result)
	if err != nil {
		t.Fatal(err)
		//t.Fatal("No error occuers, but expected error: Not found accountlogin: test2")
	}
	if result.Role != "test2" {
		t.Fatalf("Got role: %s, expected: %s", result.Role, acc.Role)
	}
	acc.Owners = append(acc.Owners, "test3")
	err = acc.Update()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Find(bson.M{"accountlogin": "test"}).One(&result)
	if err != nil {
		t.Fatal(err)
		//t.Fatal("No error occuers, but expected error: Not found accountlogin: test2")
	}
	log.Printf("TestAccount2_Update result: %+v", result)
	if result.Owners[1] != "test3" {
		t.Fatalf("Account NOT(!) updated. \n Not enough owners: Got owners: %s, expected: %s", result.Owners, "test3")
	}
}
func TestAccount2_GetInfo(t *testing.T) {
	a, err := NewAccount2("test", "test", "test", "test").GetInfo()
	if err != nil {
		t.Fatal(err)
	}
	if a.Role != "test2" {
		t.Fatalf("Unexpected role value: got %s, expected %s", a.Role, "test2")
	}
	_, err = NewAccount2("test", "test", "qwe", "test").GetInfo()
	if err == nil {
		t.Fatalf("Expected error Not found, got %v", err)
	}
	//if a.Role != "test2" {
	//	t.Fatalf("Unexpected role value: got %s, expected %s", a.Role, "test2")
	//}
}
func TestAccount2_GetAccountList(t *testing.T) {
	a := NewAccount2("qwe123test", "test", "test1", "test1")
	a.Owners = append(a.Owners, "test1")
	a.Update()
	if err != nil {
		t.Fatal(err)
	}
	a2 := NewAccount2("qwe1sadxxxst", "test", "test12", "test12")
	if err != nil {
		t.Fatal(err)
	}
	a2.Owners = append(a2.Owners, "test1")
	a2.Update()
	if err != nil {
		t.Fatal(err)
	}
	find := NewAccount2("s", "s", "s", "s")
	find.Owners = append(find.Owners, "test1")
	res, err := find.GetAccountList()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("TestAccount2_GetAccountList result list: %+v", res)
	if res[0].Creator != a.Creator {
		t.Fatalf("Not equal: %+v \n and %+v", res[0].Creator, a.Creator)
	}
	if res[1].Creator != a2.Creator {
		t.Fatalf("Not equal: %+v \n and %+v", res[1].Creator, a2.Creator)
	}
}
func TestAccount2_Remove(t *testing.T) {
	a := NewAccount2("qwe123test", "test", "test1", "test1")
	a.Owners = append(a.Owners, "test1")
	err := a.Update()
	if err != nil {
		t.Fatal(err)
	}
	res, err := a.GetAccountList()
	if err != nil {
		t.Fatal(err)
	}
	//log.Printf("TestAccount2_Remove BEFORE remove: \n %+v \n", res)
	err = a.Remove()
	if err != nil {
		t.Fatal(err)
	}
	res1, err := a.GetAccountList()
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(res1, res) {
		t.Error("Results from base are equal, but they should not.")
	}

	//logrus.Info("Using CleanDB")
	//cleanDB <- true
	//cleanUpTest()
}

//func TestcleanUpTest(t *testing.T) {
//	cleanUpTest()
//}
func cleanUpTest() {
	testDataBase := mgo.Database{
		Name:    testDB,
		Session: mainSession,
	}

	err := testDataBase.DropDatabase()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Test DB ", testDB, " droped.")
	time.Sleep(10 * time.Millisecond)
	done <- true
	return
}
