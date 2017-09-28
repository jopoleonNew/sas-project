package updaters

import (
	"log"

	"time"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/model"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

func Start(timeStep time.Duration) {
	go func(t time.Duration) {
		for {
			time.Sleep(t)
			log.Println(" <---- START YANDEX UPDATER ")
			updater()
			log.Println(" <---- END YANDEX UPDATER ")
		}
	}(timeStep)
}
func updater() {

	db := model.ImportDB

	s := db.Session.Clone()
	defer s.Close()
	c := s.DB(db.Name).C("accountsList")
	allAccounts := []model.Account2{}
	err := c.Find(nil).All(&allAccounts)
	if err != nil {
		log.Println("updater c.Find(bson.M{}).One(&allAccounts) error: ", err)
	}
	for _, acc := range allAccounts {
		go func(a model.Account2) {
			var idslice []int
			for _, id := range a.CampaignsInfo {
				idslice = append(idslice, id.ID)
			}
			account := yad.NewAccount(a.Accountlogin, a.AuthToken)
			startTime := time.Now()
			endTime := startTime.AddDate(0, 0, -1)
			statres, err := account.GetStatisticsConc(idslice, endTime, startTime)
			if err != nil {
				logrus.Errorf("updater() account.GetStatisticsConc %v error: %v", account, err)
				//return info, fmt.Errorf("CollectAccountandAddtoBD account.GetStatisticsConc %v error: %v", account, err)
			}
			err = model.UpdateYandexStatistic(a.Accountlogin, statres)
			if err != nil {
				logrus.Errorf("updater() UpdateYandexStatistic %v error: %v", account, err)
				//return info, fmt.Errorf("CollectAccountandAddtoBD SaveYandexStatistic %v error: %v", account, err)
			}

		}(acc)
	}
}
