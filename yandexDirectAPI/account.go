package yad

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Account struct {
	Login      string
	OAuthToken string `json:"token"`
}

const YandexStatThreshold int = 900

func NewAccount() *Account {
	return &Account{
		Login:      "",
		OAuthToken: "",
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

func (a *Account) GetAgencyLogins() ([]Client, error) {
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

	client := &http.Client{}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqbytes))
	if err != nil {
		return nil, errors.New("http.NewRequest error: " + err.Error())
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept-Language", "ru")
	r.Header.Add("Client-Login", a.Login)
	r.Header.Add("Authorization", "Bearer "+a.OAuthToken)
	r.Header.Add("Client-ID", application.ID)
	log.Println("makeV5GetRequest body Request: ", r)
	resp, err := client.Do(r)
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

	client := http.Client{}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqbytes))
	if err != nil {
		return nil, errors.New("GetStatistics http.NewRequest error: " + err.Error())
	}
	r.Header.Add("Content-Type", "application/json; charset=utf-8")
	//log.Println("GetStatistics body Request: ", r)
	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.New("GetStatistics http.Request.Do error: " + err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("GetStatistics ioutil.ReadAll error: " + err.Error())
	}

	result := new(ResultV4CampStat)
	err = json.Unmarshal(body, &result)
	//log.Println("GetStatistics body response: ", string(body))
	if err != nil {
		return nil, errors.New("GetStatistics json.Unmarshal error: " + err.Error())
	}
	log.Printf("GetStatistics body response: %+v", result)

	return result.Data, nil
}

func (a *Account) GetStatisticsConc(ids []int, start, end time.Time) ([]CampaignStat, error) {
	url := apiURL + "/v4/json"
	fomrmatstart := start.Format(ctLayout)
	formatend := end.Format(ctLayout)
	StartGetStatisticsConc := time.Now()
	// Checking yandex conditions:
	// The number of strings in the method response should not exceed 1000
	// Number of Campaigns IDs * Days in statistic request < YandexStatThreshold = ~1000
	delta := end.Sub(start)
	days := delta.Hours() / 24
	log.Println(days)
	//respChan := make(chan *ResultV4CampStat)
	log.Println("The number of string: ", int(days)*len(ids))
	log.Println("The fomrmatstart ", fomrmatstart)
	log.Println("The formatend ", formatend)
	log.Println("The ids slice ", ids)
	//numStrings := int(days) * len(ids)
	if int(days)*len(ids) <= 0 {
		return nil, errors.New("GetStatisticsConc The number of string is zero")
	}
	// check do we need to split int(days)*len(ids)) to parts less then YandexStatThreshold
	if int(days)*len(ids) >= YandexStatThreshold {
		log.Println("Recieved amount of days: ", int(days), " and length of campaign ID slice: ", len(ids))
		log.Println("The number of string: ", int(days)*len(ids))
		sort.Ints(ids)
		//log.Printf("%#v", ids)
		var BigErr error
		var NewIds []int
		var statsslice []CampaignStat
		// super advanced splitting algorithm
		idThreshold := (YandexStatThreshold / int(days))
		itersAmount := (len(ids) / idThreshold) + 1
		wg := sync.WaitGroup{}
		//wg.Add(itersAmount)
		// make buffered chanel to control amount of simultaneous goroutines
		sema := make(chan struct{}, 5)

		for j := 0; j < itersAmount; j++ {
			var reqIds []int
			for i := idThreshold * j; i < idThreshold*(j+1) && i < (len(ids)); i++ {
				reqIds = append(reqIds, ids[i])
				//fmt.Print("Id:=", ids[i], ", i:=", i)
				NewIds = append(NewIds, ids[i])
			}

			// because of YandexDirectAPI limitation on simultaneous connections
			// make add only 5 goroutines:
			// Технические ограничения
			// Допускается не более пяти (5) одновременных запросов к API от лица одного пользователя.

			wg.Add(1)
			go func() {
				defer wg.Done()
				sema <- struct{}{}
				defer func() { <-sema }()
				log.Println("\n", reqIds, "THE J VALUE >>>>>>>>>>", j)
				req := &RequestV4{
					Method: "GetSummaryStat",
					Token:  a.OAuthToken,
					ParamV4: &ParamV4{
						CampaignIDS: reqIds,
						StartDate:   fomrmatstart,
						EndDate:     formatend,
					},
				}

				reqbytes, err := json.Marshal(req)
				if err != nil {
					BigErr = errors.New("GetStatisticsConc json.Marshal error: " + err.Error())
					return
				}

				client := http.Client{}
				r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqbytes))
				if err != nil {
					BigErr = errors.New("GetStatisticsConc http.NewRequest error: " + err.Error())
					return
				}
				r.Header.Add("Content-Type", "application/json; charset=utf-8")

				//log.Println("GetStatistics body Request: ", r)

				sentreq := time.Now()
				resp, err := client.Do(r)
				if err != nil {
					BigErr = errors.New("GetStatisticsConc http.Request.Do error: " + err.Error())
					return
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					BigErr = errors.New("GetStatisticsConc ioutil.ReadAll error: " + err.Error())
					return
				}
				if string(body) == "" {
					BigErr = errors.New("GetStatisticsConc body response is empty")
					return
				}
				if strings.Contains(string(body), "error") {
					var errresult = YandexV4Error{}
					err := json.Unmarshal(body, &errresult)
					if err != nil {
						BigErr = errors.New("GetStatisticsConc json.Unmarshal error:  " + err.Error())
						return
					}
					//ErrorCode, _ := strconv.Atoi(errresult.ErrorCode)
					BigErr = errors.New("GetStatistics Yandex Direct API error: " + strconv.Itoa(errresult.ErrorCode) + " " + errresult.ErrorDescription + errresult.ErrorDetail)
					return
				}
				getreq := time.Since(sentreq)

				log.Println("\n\n  TIME YANDEX REQEST : ", getreq.Seconds())

				result := new(ResultV4CampStat)
				err = json.Unmarshal(body, &result)
				if err != nil {
					BigErr = errors.New("GetStatistics json.Unmarshal error: " + err.Error())
					return
				}
				statsslice = append(statsslice, result.Data...)
			}()
		}
		wg.Wait()
		//close(respChan)
		if BigErr != nil {
			return nil, BigErr
		}
		log.Println("Returning statsslice time (StartGetStatisticsConc):\n ", time.Since(StartGetStatisticsConc))

		return statsslice, nil
		//return nil, errors.New(">>>>>>>>>The number of strings in the method response should not exceed 1000")

	}

	req := &RequestV4{
		Method: "GetSummaryStat",
		Token:  a.OAuthToken,
		ParamV4: &ParamV4{
			CampaignIDS: ids,
			StartDate:   fomrmatstart,
			EndDate:     formatend,
		},
	}

	reqbytes, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("GetStatisticsConc json.Marshal error: " + err.Error())
	}

	client := http.Client{}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqbytes))
	if err != nil {
		return nil, errors.New("GetStatisticsConc http.NewRequest error: " + err.Error())
	}
	r.Header.Add("Content-Type", "application/json; charset=utf-8")
	//log.Println("GetStatisticsConc body Request: ", r)
	//sentreq := time.Now()
	//log.Println("GetStatisticsConc client.Do request body ", r)
	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.New("GetStatisticsConc http.Request.Do error: " + err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("GetStatisticsConc ioutil.ReadAll error: " + err.Error())
	}
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
		return nil, errors.New("GetStatistics Yandex Direct API error: " + strconv.Itoa(errresult.ErrorCode) + " " + errresult.ErrorDescription + errresult.ErrorDetail)
	}
	//log.Println("GetStatisticsConc client.Do body response ", string(body))

	//getreq := sentreq.Sub(sentreq)
	//log.Println("\n\n$$$$  TIME YANDEX REQEST : ", time.Since(sentreq))
	log.Println("Returning statsslice time without concurrenc :\n ", time.Since(StartGetStatisticsConc))
	result := new(ResultV4CampStat)
	err = json.Unmarshal(body, &result)
	//log.Println("GetStatistics body response: ", string(body))
	if err != nil {
		return nil, errors.New("GetStatisticsConc json.Unmarshal error: " + err.Error())
	}

	//log.Println("statsslice without concurrenc ", result)

	//log.Printf("GetStatistics body response: %+v", result)

	return result.Data, nil
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
