package model

import (
	"log"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gogs.itcloud.pro/SAS-project/sas/vkontakteAPI"
	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
	"gopkg.in/mgo.v2/bson"
)

type YandexTime struct {
	Time time.Time
}

const ctLayout = "2006-01-02"

func (ct *YandexTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	ct.Time, err = time.Parse(ctLayout, string(b))
	if err != nil {
		logrus.Fatal("models YandexTime UnmarshalJSON time.Parse error: ", err)
		return err
	}
	return nil
}

func (ct *YandexTime) MarshalJSON() ([]byte, error) {
	// if ct.Time.UnixNano() == nilTime {
	// 	return []byte("null"), nil
	// }
	return []byte(strconv.Quote(ct.Time.Format(ctLayout))), nil
}

type StatisticDataType struct {
	SessionDepthSearch    interface{} `json:"SessionDepthSearch"`
	SumSearch             float32     `json:"SumSearch"`
	ClicksContext         int         `json:"ClicksContext"`
	SessionDepthContext   interface{} `json:"SessionDepthContext"`
	StatDate              YandexTime  `json:"StatDate"`
	GoalCostSearch        interface{} `json:"GoalCostSearch"`
	GoalConversionContext interface{} `json:"GoalConversionContext"`
	ShowsContext          interface{} `json:"ShowsContext"`
	SumContext            interface{} `json:"SumContext"`
	GoalConversionSearch  interface{} `json:"GoalConversionSearch"`
	ShowsSearch           interface{} `json:"ShowsSearch"`
	CampaignID            int         `json:"CampaignID"`
	GoalCostContext       interface{} `json:"GoalCostContext"`
	ClicksSearch          int         `json:"ClicksSearch"`
}

type GetSummaryStatRes struct {
	Data []StatisticDataType `json:"Data"`
}

//implementing Sort.sort interface for GetSummaryStatRes struct
func (p GetSummaryStatRes) Len() int {
	return len(p.Data)
}

func (p GetSummaryStatRes) Less(i, j int) bool {
	return p.Data[i].StatDate.Time.Before(p.Data[j].StatDate.Time)
}

func (p GetSummaryStatRes) Swap(i, j int) {
	p.Data[i], p.Data[j] = p.Data[j], p.Data[i]
}
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
	_, err := c.Upsert(bson.M{}, struct {
		Result []yad.CampaignStat `json:"result"`
	}{stats})
	if err != nil {
		log.Println("SaveYandexStatistic input.Insert error: ", err)
		return err
	}
	return nil
}

func UpdateYandexStatistic(accountlogin string, stats []yad.CampaignStat) error {
	s := mainSession.Clone()
	defer s.Close()
	colname := "yandex" + accountlogin + "stats"
	c := s.DB("statistic").C(colname)
	//c.Find
	//colQuerier := bson.M{"accountlogin": a.Accountlogin, "source": a.Source}
	change1 := bson.M{"$addToSet": bson.M{"result": stats}}
	//bson.M{"$addToSet": bson.M{"owners": bson.M{"$each": a.Owners}}}
	_, err := c.Upsert(bson.M{}, change1)
	//err := c.Update("result", change1)
	if err != nil {
		logrus.Errorf("UpdateYandexStatistic c.Upsert('result', change1) err: %+v", err)
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
