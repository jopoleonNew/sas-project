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

type CampaignStat struct {
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
	Data []CampaignStat `json:"Data"`
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

type Statistic interface {
	Save() error
	Get() error
	Update() error
}

type VKstatistic struct {
	Accountlogin string
	Stats        []vkontakteAPI.AdStatistic
}
type Yandexstatistic struct {
	Accountlogin string
	Stats        []yad.CampaignStat
}

func (vk *VKstatistic) Save() {

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

	var save = make([]interface{}, len(stats))
	for i, s := range stats {
		save[i] = s
	}
	//err := c.Insert(struct {
	//	Result []yad.CampaignStat `json:"result"`
	//}{stats})
	//if err != nil {
	//	log.Println("SaveYandexStatistic input.Insert error: ", err)
	//	return err
	//}
	err = c.Insert(save...)
	if err != nil {
		log.Println("SaveYandexStatistic  c.Insert([]interface{})error: ", err)
		return err
	}
	//var inf []interface{}
	//byte, err := json.Marshal(stats)
	//if err != nil {
	//	log.Println("SaveYandexStatistic input.Insert error: ", err)
	//	return err
	//}
	//err = json.Unmarshal((byte), &inf)
	//if err != nil {
	//	log.Println("SaveYandexStatistic input.Insert error: ", err)
	//	return err
	//}
	//err = c.Insert(struct {
	//	Result []yad.CampaignStat `json:"result"`
	//}{stats})

	return nil
}

func UpdateYandexStatistic(accountlogin string, stats []yad.CampaignStat) error {
	s := mainSession.Clone()
	defer s.Close()
	colname := "yandex" + accountlogin + "stats"
	c := s.DB("statistic").C(colname)
	for _, s := range stats {
		_, err = c.UpdateAll(bson.M{"campaignid": s.CampaignID, "statdate": bson.M{"time": s.StatDate.Time}}, bson.M{"$set": s})
		//err := c.Update("result", change1)
		if err != nil {
			logrus.Errorf("UpdateYandexStatistic c.UpdateAll error: %+v, \n for state : %+v", err, s)
			return err
		}
	}
	return nil
}

func GetYandexStatistic(accountlogin string) ([]yad.CampaignStat, error) {
	var out []yad.CampaignStat
	s := mainSession.Clone()
	defer s.Close()
	colname := "yandex" + accountlogin + "stats"
	c := s.DB("statistic").C(colname)
	err := c.Find(nil).All(&out)
	if err != nil {
		log.Println("GetYandexStatistic input.Insert error: ", err)
		return out, err
	}
	return out, nil
}
