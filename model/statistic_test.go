package model

import (
	"testing"

	"reflect"

	"time"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

func TestSaveYandexStatistic(t *testing.T) {
	test := []yad.CampaignStat{
		{CampaignID: 123,
			ClicksContext: 123,
			ClicksSearch:  123,
			StatDate:      yad.YandexTime{time.Date(2017, time.September, 20, 20, 00, 00, 00, time.FixedZone("test", 0))},
			SumSearch:     123},
		{CampaignID: 123,
			ClicksContext: 13,
			StatDate:      yad.YandexTime{time.Date(2017, time.September, 21, 20, 00, 00, 00, time.FixedZone("test", 0))},
			SumSearch:     13},
		{CampaignID: 1323,
			ShowsSearch: 1723,
			StatDate:    yad.YandexTime{time.Date(2017, time.September, 20, 20, 00, 00, 00, time.FixedZone("test", 0))},
			SumContext:  17623,
			SumSearch:   1283},
		{CampaignID: 1323,
			GoalCostSearch:      2,
			SessionDepthContext: 2,
			SessionDepthSearch:  2,
			StatDate:            yad.YandexTime{time.Date(2017, time.September, 20, 20, 00, 00, 00, time.FixedZone("test", 0))},
		},
	}
	testAccLogin := "testSave"
	err := SaveYandexStatistic(testAccLogin, test)
	if err != nil {
		t.Fatal("SaveYandexStatistic error:", err)
	}
	var out struct {
		Result []yad.CampaignStat `json:"result"`
	}
	s := mainSession.Clone()
	defer s.Close()
	colname := "yandex" + testAccLogin + "stats"
	c := s.DB("statistic").C(colname)
	err = c.Find(nil).One(&out)
	if err != nil {
		t.Fatal("c.Find(nil).One(&out) error:", err)
	}
	if !reflect.DeepEqual(test, out.Result) {
		t.Errorf("SaveYandexStatistic Statistic not equal:\n got %+v, \n expected %+v", out.Result, test)
	}
	//c.DropCollection()
}
func TestGetYandexStatistic(t *testing.T) {
	test := []yad.CampaignStat{
		{CampaignID: 123,
			ClicksContext: 123,
			ClicksSearch:  123,
			StatDate:      yad.YandexTime{time.Date(2017, time.September, 20, 20, 00, 00, 00, time.FixedZone("test", 0))},
			SumSearch:     123},
		{CampaignID: 123,
			ClicksContext: 13,
			StatDate:      yad.YandexTime{time.Date(2017, time.September, 21, 20, 00, 00, 00, time.FixedZone("test", 0))},
			SumSearch:     13},
		{CampaignID: 1323,
			ShowsSearch: 1723,
			StatDate:    yad.YandexTime{time.Date(2017, time.September, 20, 20, 00, 00, 00, time.FixedZone("test", 0))},
			SumContext:  17623,
			SumSearch:   1283},
		{CampaignID: 1323,
			GoalCostSearch:      2,
			SessionDepthContext: 2,
			SessionDepthSearch:  2,
			StatDate:            yad.YandexTime{time.Date(2017, time.September, 20, 20, 00, 00, 00, time.FixedZone("test", 0))},
		},
	}
	testAccLogin := "testSave"
	res, err := GetYandexStatistic(testAccLogin)
	if err != nil {
		t.Fatal("GetYandexStatistic(testAccLogin) error:", err)
	}
	if !reflect.DeepEqual(test, res) {
		t.Errorf("GetYandexStatistic result not correct:\n got %+v, \n expected %+v", res, test)
	}
}

func TestUpdateYandexStatistic(t *testing.T) {
	testUpdate := []yad.CampaignStat{
		{CampaignID: 123,
			ClicksContext:         123,
			ClicksSearch:          123,
			GoalConversionContext: 123,
			GoalConversionSearch:  123,
			GoalCostContext:       123,
			GoalCostSearch:        123,
			SessionDepthContext:   123,
			SessionDepthSearch:    123,
			ShowsContext:          123,
			ShowsSearch:           123,
			StatDate:              yad.YandexTime{time.Date(2017, time.September, 20, 20, 00, 00, 00, time.FixedZone("test", 0))},
			SumSearch:             123,
			SumContext:            123,
		}}
	testAccLogin := "testSave"
	beforeUpdate, err := GetYandexStatistic(testAccLogin)
	if err != nil {
		t.Fatal("GetYandexStatistic(testAccLogin) error:", err)
	}
	err = UpdateYandexStatistic(testAccLogin, testUpdate)
	if err != nil {
		t.Fatal("SaveYandexStatistic error:", err)
	}
	afterUpdate, err := GetYandexStatistic(testAccLogin)
	if err != nil {
		t.Fatal("GetYandexStatistic(testAccLogin) error:", err)
	}
	logrus.Infof("Before : %+v \n After : %+v", beforeUpdate, afterUpdate)
	//if !reflect.DeepEqual(testUpdate, res) {
	//	t.Errorf("Statistic not equal:\n got %+v, \n expected %+v", res, testUpdate)
	//}
}

//func SaveYandexStatistic(accountlogin string, stats []yad.CampaignStat) error {
//	s := mainSession.Clone()
//	defer s.Close()
//	colname := "yandex" + accountlogin + "stats"
//	c := s.DB("statistic").C(colname)
//	_, err := c.Upsert(bson.M{}, struct {
//		Result []yad.CampaignStat `json:"result"`
//	}{stats})
//	if err != nil {
//		log.Println("SaveYandexStatistic input.Insert error: ", err)
//		return err
//	}
//	return nil
//}
//
//func UpdateYandexStatistic(accountlogin string, stats []yad.CampaignStat) error {
//	s := mainSession.Clone()
//	defer s.Close()
//	colname := "yandex" + accountlogin + "stats"
//	c := s.DB("statistic").C(colname)
//	//c.Find
//	//colQuerier := bson.M{"accountlogin": a.Accountlogin, "source": a.Source}
//	change1 := bson.M{"$addToSet": bson.M{"result": stats}}
//	_, err := c.Upsert(bson.M{}, change1)
//	//err := c.Update("result", change1)
//	if err != nil {
//		logrus.Errorf("UpdateYandexStatistic c.Upsert('result', change1) err: %+v", err)
//		return err
//	}
//	return nil
//}
