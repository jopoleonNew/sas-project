package model

import (
	"log"
	"strconv"
	//"git.itcloud.pro/egortictac/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

func MakeStatisticCollection(username, acclogin string, campaingStat GetSummaryStatRes) error {
	log.Println("MakeStatisticCollection used")

	s := mainSession.Clone()
	defer s.Close()

	var campId string
	log.Println("MakeStatisticCollection LENGTH:", len(campaingStat.Data))
	//log.Println("MakeStatisticCollection value:", campaingStat.Data)
	log.Printf("\n MakeStatisticCollection VALUE: %+v \n", campaingStat)
	if len(campaingStat.Data) > 0 {

		campId = strconv.Itoa(campaingStat.Data[0].CampaignID)
		log.Println("MakeStatisticCollection campID:", campId)
	} else {
		campId = "empty" + RandStringBytes(5)
		log.Println("MakeStatisticCollection campId = empty ", campId)
	}
	collectionname := username + acclogin + campId
	log.Println("MakeStatisticCollection collectionname", collectionname)
	c := s.DB(mainDB.Name).C(collectionname)
	err := c.Insert(campaingStat)
	if err != nil {
		log.Println("MakeStatisticCollection input.Insert error: ", err)
		return err
	}
	return nil
}

type StatisticSaver interface {
	Save()
	Get()
}

func SaveVKStatistic(accountlogin string, stats vkontakteAPI.AdStatistic) error {

	s := mainSession.Clone()
	defer s.Close()
	colname := "vk" + accountlogin + "stats"
	c := s.DB("statistic").C(colname)
	err := c.Insert(stats)
	if err != nil {
		log.Println("SaveVKStatistic input.Insert error: ", err)
		return err
	}
	return nil
}
func GetVKStatistic(accountlogin string) (vkontakteAPI.AdStatistic, error) {
	var out vkontakteAPI.AdStatistic
	s := mainSession.Clone()
	defer s.Close()
	colname := "vk" + accountlogin + "stats"
	c := s.DB("statistic").C(colname)
	err := c.Find(nil).One(&out)
	if err != nil {
		log.Println("GetVKStatistic input.Insert error: ", err)
		return out, err
	}
	return out, nil
}

func SaveYandexStatistic(accountlogin string, stats []yad.CampaignStat) error {
	s := mainSession.Clone()
	defer s.Close()
	colname := "yandex" + accountlogin + "stats"
	c := s.DB("statistic").C(colname)
	err := c.Insert(struct {
		Result []yad.CampaignStat `json:"result"`
	}{stats})
	if err != nil {
		log.Println("SaveYandexStatistic input.Insert error: ", err)
		return err
	}
	return nil
}
func GetYandexStatistic(accountlogin string) ([]yad.CampaignStat, error) {
	var out struct {
		Result []yad.CampaignStat `json:"result"`
	}
	s := mainSession.Clone()
	defer s.Close()
	colname := "yandex" + accountlogin + "stats"
	c := s.DB("statistic").C(colname)
	err := c.Find(nil).One(&out)
	if err != nil {
		log.Println("GetYandexStatistic input.Insert error: ", err)
		return out.Result, err
	}
	return out.Result, nil
}
