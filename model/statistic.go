package model

import (
	"log"
	"strconv"
	//"git.itcloud.pro/egortictac/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
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

func SaveVKStatistic(accountlogin string, stats vkontakteAPI.AdStatistic) error {

	s := mainSession.Clone()
	defer s.Close()
	colname := accountlogin + "vk" + "stats"
	c := s.DB(mainDB.Name).C(colname)
	err := c.Insert(stats)
	if err != nil {
		log.Println("MakeStatisticCollection input.Insert error: ", err)
		return err
	}
	return nil
}
