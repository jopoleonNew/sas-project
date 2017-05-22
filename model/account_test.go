package model

import (
	"log"
	"testing"

	"gogs.itcloud.pro/SAS-project/sas/app"
)

func init() {

	app.InitConf("../configuration.json")
	Config = app.GetConfig()
	log.Printf("CONFIG FILE MAIN: %+v", Config)
	err := SetDBParams(Config.Mongourl, Config.DBname)
	if err != nil {
		log.Fatal(err)
	}
}

//
//func TestAccount_IsExist(t *testing.T) {
//	acc := NewAccount()
//	acc.Username = "test"
//	acc.Accountlogin = "test"
//	acc.Source = "test"
//	acc.YandexRole = "test"
//	//err = acc.Remove()
//	//if err != nil {
//	//	t.Error(err)
//	//}
//	exist, err := acc.IsExist()
//	if err == ErrAccNotFound {
//
//	} else {
//		t.Error(err)
//	}
//	expected := false
//	if exist != expected {
//		t.Fatalf("Expected %v, got %v", expected, exist)
//	}
//	err = acc.AdvanceUpdate()
//	exist2, err := acc.IsExist()
//	if err == ErrAccNotFound {
//
//	} else {
//		t.Error(err)
//	}
//	expected2 := true
//	if exist2 != expected2 {
//		t.Fatalf("Expected %v, got %v", expected2, exist2)
//	}
//	err = acc.Remove()
//	if err != nil {
//		t.Error(err)
//	}
//}

func BenchmarkAccount_IsExist40(b *testing.B) {
	acc := NewAccount()
	acc.Username = "test"
	acc.Accountlogin = "test"
	acc.Source = "test"
	acc.YandexRole = "test"
	err = acc.AdvanceUpdate()
	if err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		_, err := acc.IsExist()
		if err != nil {
			b.Error(err)
		}
	}
	err = acc.Remove()
	if err != nil {
		b.Error(err)
	}

}
