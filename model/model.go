package model

import (
	"log"
	"strconv"
	"time"
)

type Members struct {
	Name string
	Role string
}

type TemplateInfoStruct struct {
	CurrentUser string
	UsingReport string
	AccountList []Account
	FullStats   []GetSummaryStatRes
}

type TemplateInfo struct {
	CurrentUser    string
	CurrentAccount Account2
	UsingReport    string
	AccountList    []Account2
	FullStats      []GetSummaryStatRes
}
type Campaign struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Ads    []Ad   `json:"ads"`
	Owner  string `json:"owner"`
}

type Ad struct {
	ID     int
	CampID int
	Name   string
	Status string
	//cpc - цена за переход в копейках.
	CPC int
	//cpm - цена за 1000 показов в копейках.
	CPM  int
	Type string
}

//id
//integer	идентификатор объявления.
//campaign_id
//integer	идентификатор кампании.
//ad_format
//integer	формат объявления. Возможные значения:
//1 — изображение и текст;
//2 — большое изображение;
//3 — эксклюзивный формат;
//4 — продвижение сообществ или приложений, квадратное изображение;
//5 — приложение в новостной ленте (устаревший);
//6 — мобильное приложение;
//9 — запись в сообществе.
//cost_type
//integer, [0,1]	тип оплаты. Возможные значения:
//0 — оплата за переходы;
//1 — оплата за показы.
//cpc
//integer	(если cost_type = 0) цена за переход в копейках.
//cpm
//integer	(если cost_type = 1) цена за 1000 показов в копейках.
//impressions_limit
//integer	(если задано) ограничение количества показов данного объявления на одного пользователя. Может присутствовать для некоторых форматов объявлений, для которых разрешена установка точного значения.
//impressions_limited
//integer, [1]	(если задано) признак того, что количество показов объявления на одного пользователя ограничено. Может присутствовать для некоторых объявлений, для которых разрешена установка ограничения, но не разрешена установка точного значения. 1 — не более 100 показов на одного пользователя.
//ad_platform	(если значение применимо к данному формату объявления) рекламные площадки, на которых будет показываться объявление. Возможные значения:
//(если ad_format равен 1)
//0 — ВКонтакте и сайты-партнёры;
//1 — только ВКонтакте.
//(если ad_format равен 9)
//all — все площадки;
//desktop — полная версия сайта;
//mobile — мобильный сайт и приложения.
//ad_platform_no_wall
//integer, [1]	1 — для объявления задано ограничение «Не показывать на стенах сообществ».
//all_limit
//integer	общий лимит объявления в рублях. 0 — лимит не задан.
//category1_id
//integer	ID тематики или подраздела тематики объявления. См. ads.getCategories.
//category2_id
//integer	ID тематики или подраздела тематики объявления. Дополнительная тематика.
//status
//integer	статус объявления. Возможные значения:
//0 — объявление остановлено;
//1 — объявление запущено;
//2 — объявление удалено.
//name
//string	название объявления.
//approved
//integer	статус модерации объявления. Возможные значения:
//0 — объявление не проходило модерацию;
//1 — объявление ожидает модерации;
//2 — объявление одобрено;
//3 — объявление отклонено.
//video
//integer, [1]	1 — объявление является видеорекламой.
//disclaimer_medical
//integer, [1]	1 — включено отображение предупреждения:
//«Есть противопоказания. Требуется консультация специалиста.»
//disclaimer_specialist
//integer, [1]	1 — включено отображение предупреждения:
//«Необходима консультация специалистов.»
//disclaimer_supplements
//integer, [1]	1 — включено отображение предупреждения:
//«БАД. Не является лекарственным препаратом.»

// YandexTokenbody is used in MakeYandexOauthRequest()
// to unmarshal yandex response body and get AccessToken

// "CampaignID": (int),
//        "StatDate": (date),
//        "SumSearch": (float),
//        "SumContext": (float),
//        "ShowsSearch": (int),
//        "ShowsContext": (int),
//        "ClicksSearch": (int),
//        "ClicksContext": (int),
//        "SessionDepthSearch": (float),
//        "SessionDepthContext": (float),
//        "GoalConversionSearch": (float),
//        "GoalConversionContext": (float),
//        "GoalCostSearch": (float),
//        "GoalCostContext": (float)
//  SessionDepthSearch    string  `json:"SessionDepthSearch"`
// SumSearch             int     `json:"SumSearch"`
// ClicksContext         int     `json:"ClicksContext"`
// SessionDepthContext   string  `json:"SessionDepthContext"`
// StatDate              string  `json:"StatDate"`
// GoalCostSearch        float32 `json:"GoalCostSearch"`
// GoalConversionContext string  `json:"GoalConversionContext"`
// ShowsContext          int     `json:"ShowsContext"`
// SumContext            int     `json:"SumContext"`
// GoalConversionSearch  string  `json:"GoalConversionSearch"`
// ShowsSearch           int     `json:"ShowsSearch"`
// CampaignID            int     `json:"CampaignID"`
// GoalCostContext       float32 `json:"GoalCostContext"`
// ClicksSearch          int     `json:"ClicksSearch"`
// interface{}

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
		log.Fatal("models UnmarshalJSON time.Parse error: ", err)
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
