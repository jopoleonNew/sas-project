package yad

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Account struct {
	Login      string
	OAuthToken string `json:"token"`
}

const YandexStatThreshold int = 900

func NewAccount(login, oauthtoken string) *Account {
	return &Account{
		Login:      login,
		OAuthToken: oauthtoken,
	}
}

func (a *Account) GetCampaignList() ([]Campaign, error) {
	url := apiURL + "/json/v5/campaigns"
	fieldNames := []string{"Id", "Name", "Status"}

	result := new(ResultV5)

	body, err := a.makeV5GetRequest(url, fieldNames)
	if err != nil {
		return nil, errors.New("GetCampaignList " + err.Error())
	}
	//log.Println("GetCampaignList body from a.makeV5GetRequest: ", string(body))
	if err := json.Unmarshal(body, result); err != nil {
		return nil, errors.New("Respond unmarshal error: " + err.Error())
	}
	//log.Println("GetCampaignList after Unmarshal: ", result)
	return result.Result.Campaigns, nil
}

func (a *Account) GetAgencyLogins() ([]ClientAg, error) {
	url := apiURL + "/json/v5/agencyclients"
	fieldNames := []string{"Login", "Representatives"}

	result := new(ResultV5)

	body, err := a.makeV5GetRequest(url, fieldNames)
	if err != nil {
		return nil, errors.New("GetAgencyLogins " + err.Error())
	}

	if err := json.Unmarshal(body, result); err != nil {
		log.Println("GetAgencyLogins makeV5GetRequest Response unmarshale error: ", err, " trying to reunmarshal by YandexV5Error")
		var errresult = YandexV5Error{}
		err := json.Unmarshal(body, &errresult)
		if err != nil {
			return nil, errors.New("GetAgencyLogins " + err.Error())
		}
		log.Println("GetAgencyLogins makeV5GetRequest Yandex Direct API error: " + errresult.ErrorCode + errresult.ErrorString + errresult.ErrorDescription + errresult.RequsetID)
		return nil, errors.New("Yandex Direct API error: " + errresult.ErrorCode + errresult.ErrorString + errresult.ErrorDescription + errresult.RequsetID)
	}

	return result.Result.Clients, nil
}

func (a *Account) makeV5GetRequest(url string, fieldNames []string) ([]byte, error) {
	req := &RequestV5{
		Method: "get",
		ParamsV5: &ParamsV5{
			SelectionCriteria: struct{}{},
			FieldNames:        fieldNames,
		},
	}

	reqbytes, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("json.Marshal error: " + err.Error())
	}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqbytes))
	if err != nil {
		return nil, errors.New("http.NewRequest error: " + err.Error())
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept-Language", "ru")
	r.Header.Add("Client-Login", a.Login)
	r.Header.Add("Authorization", "Bearer "+a.OAuthToken)
	r.Header.Add("Client-ID", application.ID)

	//log.Println("makeV5GetRequest body Request: ", r)
	result, err := YPool.SendWork(r)
	if err != nil {
		return []byte{}, err
	}
	resp, ok := result.(*http.Response)
	if !ok {
		return []byte{}, fmt.Errorf("result from worker is not *http.Response type, %v ", result)
	}
	//resp, err := client.Do(r)

	if err != nil {
		return nil, errors.New("http.Request.Do error: " + err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ioutil.ReadAll error: " + err.Error())
	}
	if string(body) == "" {
		return nil, errors.New("makeV5GetRequest body response is empty")
	}
	//log.Println("makeV5GetRequest Response: ", string(body))
	if strings.Contains(string(body), "error") {
		var errresult = YandexV5Error{}
		err := json.Unmarshal(body, &errresult)
		if err != nil {
			return nil, errors.New("makeV5GetRequest json.Unmarshal error:  " + err.Error())
		}
		return body, errors.New("makeV5GetRequest Yandex Direct API error: " + errresult.ErrorCode + errresult.ErrorString + errresult.ErrorDescription + errresult.RequsetID + string(body))
	}
	//log.Println("makeV5GetRequest body response: ", string(body))
	return body, nil
}

// GetStatistics makes request to YandexAPI to get statistic by given campaign's ID in range of given dates
// StartDate	Начальная дата отчетного периода, за который возвращается статистика (YYYY-MM-DD)
// EndDate	Конечная дата отчетного периода, за который возвращается статистика (YYYY-MM-DD).
func (a *Account) GetStatistics(ids []int, start, end string) ([]CampaignStat, error) {
	url := apiURL + "/v4/json"

	req := &RequestV4{
		Method: "GetSummaryStat",
		Token:  a.OAuthToken,
		ParamV4: &ParamV4{
			CampaignIDS: ids,
			StartDate:   start,
			EndDate:     end,
		},
	}
	reqbytes, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("GetStatistics json.Marshal error: " + err.Error())
	}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqbytes))
	if err != nil {
		return nil, errors.New("GetStatistics http.NewRequest error: " + err.Error())
	}
	r.Header.Add("Content-Type", "application/json; charset=utf-8")
	resultWork, err := YPool.SendWork(r)
	if err != nil {
		return nil, err
	}
	resp, ok := resultWork.(*http.Response)
	if !ok {
		return nil, fmt.Errorf("result from worker is not *http.Response type, %v ", resultWork)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("GetStatistics ioutil.ReadAll error: " + err.Error())
	}
	//logrus.Info("GetStatistics response body: ", string(body))
	if string(body) == "" {
		return nil, errors.New("GetStatisticsConc body response is empty")
	}
	if strings.Contains(string(body), "error") {
		var errresult = YandexV4Error{}
		err := json.Unmarshal(body, &errresult)
		if err != nil {
			return nil, errors.New("GetStatisticsConc json.Unmarshal error:  " + err.Error())
		}
		//ErrorCode, _ := strconv.Atoi(errresult.ErrorCode)
		return nil, errors.New("GetStatistics Yandex Direct API error: " + strconv.Itoa(errresult.ErrorCode) + " " + errresult.ErrorDescription + " " + errresult.ErrorDetail)
	}
	result := new(ResultV4CampStat)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.New("GetStatistics json.Unmarshal error: " + err.Error())
	}

	return result.Data, nil
}

func (a *Account) GetStatisticsConc(ids []int, start, end time.Time) ([]CampaignStat, error) {
	fomrmatstart := start.Format(ctLayout)
	formatend := end.Format(ctLayout)

	// Checking yandex conditions:
	// The number of strings in the method response should not exceed 1000
	// Number of Campaigns IDs * Days in statistic request < YandexStatThreshold = ~1000
	delta := end.Sub(start)
	days := delta.Hours() / 24
	//	log.Println(days)
	//respChan := make(chan *ResultV4CampStat)
	//log.Println("The number of strings: ", int(days)*len(ids))
	//log.Println("The fomrmatstart ", fomrmatstart)
	//log.Println("The formatend ", formatend)
	//log.Println("The ids slice ", ids)
	//numStrings := int(days) * len(ids)
	if int(days)*len(ids) <= 0 {
		return nil, errors.New("GetStatisticsConc The number of strings is zero")
	}
	// check do we need to split int(days)*len(ids)) to parts less then YandexStatThreshold
	if int(days)*len(ids) >= YandexStatThreshold {
		log.Println("Recieved amount of days: ", int(days), " and length of campaign ID slice: ", len(ids))
		log.Println("The number of string: ", int(days)*len(ids))
		sort.Ints(ids)
		//log.Printf("%#v", ids)

		var NewIds []int
		var statsslice []CampaignStat
		resultChan := make(chan []CampaignStat)
		// super advanced splitting algorithm
		idThreshold := (YandexStatThreshold / int(days))
		// itersAmount - number of iterations to cover all campaings ids
		itersAmount := (len(ids) / idThreshold) + 1
		wg := sync.WaitGroup{}
		//wg.Add(itersAmount)
		// make buffered chanel to control amount of simultaneous goroutines
		//sema := make(chan struct{}, 5)

		// starting goroutine to listen result channel (resultChan) for appending full
		// account statistic for later return as result
		go func() {
			for res := range resultChan {
				statsslice = append(statsslice, res...)
			}
		}()
		for j := 0; j < itersAmount; j++ {
			var reqIds []int
			for i := idThreshold * j; i < idThreshold*(j+1) && i < (len(ids)); i++ {
				reqIds = append(reqIds, ids[i])
				//fmt.Print("Id:=", ids[i], ", i:=", i)
				NewIds = append(NewIds, ids[i])
			}

			wg.Add(1)
			go func(ids []int, start, end string) {
				if len(ids) != 0 {
					resultStatistic, err := a.GetStatistics(reqIds, start, end)
					if err != nil {
						logrus.Error("GetStatisticsConc GetStatistics error: ", err)
						return
					}
					resultChan <- resultStatistic
					//statsslice = append(statsslice, resultStatistic...)
				} else {
					logrus.Warn("a.GetStatistics(reqIds) reqIds is 0 length: ", reqIds)
					resultChan <- []CampaignStat{}
					//statsslice = append(statsslice, []CampaignStat{}...)
				}
				wg.Done()
			}(reqIds, fomrmatstart, formatend)
		}

		wg.Wait()
		close(resultChan)
		return statsslice, nil

	}

	resultStatistic, err := a.GetStatistics(ids, fomrmatstart, formatend)
	if err != nil {
		return nil, errors.New("GetStatisticsConc a.GetStatistics error: " + err.Error())
	}
	return resultStatistic, nil
}

func checkEq(a, b []int) bool {

	if a == nil && b == nil {
		log.Println("a == nil && b == nil")
		return true
	}

	if a == nil || b == nil {
		log.Println("a == nil || b == nil")
		return false
	}

	if len(a) != len(b) {
		log.Println(len(a), len(b))
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			log.Println(a[i], b[i])
			return false
		}
	}

	return true
}
