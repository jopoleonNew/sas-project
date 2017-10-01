package model

import (
	"testing"

	"time"

	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

func TestSaveYandexStatistic(t *testing.T) {
	test := []yad.CampaignStat{
		{CampaignID: 123,
			ClicksContext:      123,
			ClicksSearch:       123,
			SessionDepthSearch: 2,
			StatDate:           yad.YandexTime{time.Date(2017, time.September, 23, 00, 00, 00, 00, time.FixedZone("UTC", 0000))},
			SumSearch:          123},
		{CampaignID: 123,
			ClicksContext: 13,
			//StatDate:      yad.YandexTime{},
			SumSearch: 13},
		{CampaignID: 1323,
			ShowsSearch:        1723,
			SessionDepthSearch: 2,
			StatDate:           yad.YandexTime{time.Date(2017, time.September, 20, 00, 00, 00, 00, time.FixedZone("UTC", 0000))},
			SumContext:         17623,
			SumSearch:          1283},
		{CampaignID: 1323,
			GoalCostSearch:      2,
			SessionDepthContext: 2,
			SessionDepthSearch:  2,
			StatDate:            yad.YandexTime{time.Date(2017, time.September, 21, 00, 00, 00, 00, time.FixedZone("UTC", 0000))},
		},
	}
	testAccLogin := "testSave"
	err := SaveYandexStatistic(testAccLogin, test)
	if err != nil {
		t.Fatal("SaveYandexStatistic error:", err)
	}
	var out []yad.CampaignStat
	s := mainSession.Clone()
	defer s.Close()
	colname := "yandex" + testAccLogin + "stats"
	c := s.DB("statistic").C(colname)
	err = c.Find(nil).All(&out)
	if err != nil {
		t.Fatal("c.Find(nil).One(&out) error:", err)
	}
	for i := range test {
		if test[i].SessionDepthSearch == out[i].SessionDepthSearch &&
			test[i].SumSearch == out[i].SumSearch &&
			test[i].ClicksContext == out[i].ClicksContext &&
			test[i].SessionDepthContext == out[i].SessionDepthContext &&
			test[i].GoalCostSearch == out[i].GoalCostSearch &&
			test[i].GoalConversionContext == out[i].GoalConversionContext &&
			test[i].StatDate.Time.Unix() == out[i].StatDate.Time.Unix() &&
			test[i].ShowsContext == out[i].ShowsContext &&
			test[i].SumContext == out[i].SumContext &&
			test[i].GoalConversionSearch == out[i].GoalConversionSearch &&
			test[i].ShowsSearch == out[i].ShowsSearch &&
			test[i].CampaignID == out[i].CampaignID &&
			test[i].GoalCostContext == out[i].GoalCostContext &&
			test[i].ClicksSearch == out[i].ClicksSearch {

		} else {
			t.Errorf("SaveYandexStatistic Statistic not equal:\n got_____ %+v, \n expected %+v", out, test)
		}
	}
}

func TestGetYandexStatistic(t *testing.T) {
	test := []yad.CampaignStat{
		{CampaignID: 123,
			ClicksContext:      123,
			ClicksSearch:       123,
			SessionDepthSearch: 2,
			StatDate:           yad.YandexTime{time.Date(2017, time.September, 23, 00, 00, 00, 00, time.FixedZone("UTC", 0000))},
			SumSearch:          123},
		{CampaignID: 123,
			ClicksContext: 13,
			//StatDate:      yad.YandexTime{},
			SumSearch: 13},
		{CampaignID: 1323,
			ShowsSearch:        1723,
			SessionDepthSearch: 2,
			StatDate:           yad.YandexTime{time.Date(2017, time.September, 20, 00, 00, 00, 00, time.FixedZone("UTC", 0000))},
			SumContext:         17623,
			SumSearch:          1283},
		{CampaignID: 1323,
			GoalCostSearch:      2,
			SessionDepthContext: 2,
			SessionDepthSearch:  2,
			StatDate:            yad.YandexTime{time.Date(2017, time.September, 21, 00, 00, 00, 00, time.FixedZone("UTC", 0000))},
		},
	}
	testAccLogin := "testSave"
	res, err := GetYandexStatistic(testAccLogin)
	if err != nil {
		t.Fatal("GetYandexStatistic(testAccLogin) error:", err)
	}
	for i := range test {
		if test[i].SessionDepthSearch == res[i].SessionDepthSearch &&
			test[i].SumSearch == res[i].SumSearch &&
			test[i].ClicksContext == res[i].ClicksContext &&
			test[i].SessionDepthContext == res[i].SessionDepthContext &&
			test[i].GoalCostSearch == res[i].GoalCostSearch &&
			test[i].GoalConversionContext == res[i].GoalConversionContext &&
			test[i].StatDate.Time.Unix() == res[i].StatDate.Time.Unix() &&
			test[i].ShowsContext == res[i].ShowsContext &&
			test[i].SumContext == res[i].SumContext &&
			test[i].GoalConversionSearch == res[i].GoalConversionSearch &&
			test[i].ShowsSearch == res[i].ShowsSearch &&
			test[i].CampaignID == res[i].CampaignID &&
			test[i].GoalCostContext == res[i].GoalCostContext &&
			test[i].ClicksSearch == res[i].ClicksSearch {

		} else {
			t.Errorf("SaveYandexStatistic Statistic not equal:\n got_____ %+v, \n expected %+v", res, test)
		}

	}
}

func TestUpdateYandexStatistic(t *testing.T) {
	testUpdate := []yad.CampaignStat{
		{CampaignID: 123,
			ClicksContext:         123,
			ClicksSearch:          9999,
			GoalConversionContext: 9999,
			GoalConversionSearch:  9999,
			GoalCostContext:       9999,
			GoalCostSearch:        9999,
			SessionDepthContext:   9999,
			SessionDepthSearch:    9999,
			ShowsContext:          99999,
			ShowsSearch:           99292937743,
			StatDate:              yad.YandexTime{time.Date(2017, time.September, 23, 00, 00, 00, 00, time.FixedZone("UTC", 0000))},
			SumSearch:             99999,
			SumContext:            99999,
		}}
	testAccLogin := "testSave"
	err = UpdateYandexStatistic(testAccLogin, testUpdate)
	if err != nil {
		t.Fatal("SaveYandexStatistic error:", err)
	}
	afterUpdate, err := GetYandexStatistic(testAccLogin)
	if err != nil {
		t.Fatal("GetYandexStatistic(testAccLogin) error:", err)
	}
	var updatedFound = false
	for _, s := range afterUpdate {
		if s.CampaignID == testUpdate[0].CampaignID &&
			s.StatDate.Time.Unix() == testUpdate[0].StatDate.Time.Unix() &&
			s.ClicksSearch == testUpdate[0].ClicksSearch &&
			s.GoalConversionContext == testUpdate[0].GoalConversionContext &&
			s.GoalConversionSearch == testUpdate[0].GoalConversionSearch &&
			s.GoalCostContext == testUpdate[0].GoalCostContext &&
			s.GoalCostSearch == testUpdate[0].GoalCostSearch {
			updatedFound = true
		}
	}
	if !updatedFound {
		t.Errorf("UpdateYandexStatistic did not update statistic:, \n got %+v, \n expected: %+v ", afterUpdate, testUpdate)
	}

	s := mainSession.Clone()
	defer s.Close()
	colname := "yandex" + testAccLogin + "stats"
	c := s.DB("statistic").C(colname)
	c.DropCollection()

}
