package agencyclients

import (
	"encoding/json"
	"errors"

	yad "gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

type Service struct {
	Client *yad.Client
}

func New(c *yad.Client) Service {
	return Service{c}
}

func (s Service) Get(p GetRequest) (GetResponse, error) {
	if r, err := s.Client.Do("agencyclients", "get", p); err != nil {
		return GetResponse{}, err
	} else {
		o := map[string]GetResponse{}
		if err := json.Unmarshal(r, &o); err != nil {
			//log.Println("....../////////Agencyclients Get Unmarshal error:", err, r)
			return GetResponse{}, err
		}
		if res, ok := o["result"]; ok {
			return res, nil
		} else {
			return GetResponse{}, errors.New("Unable to find result")
		}
	}

}

func (s Service) Add(p AddRequest) (AddResponse, error) {
	if r, err := s.Client.Do("agencyclients", "add", p); err != nil {
		return AddResponse{}, err
	} else {
		o := map[string]AddResponse{}
		if err := json.Unmarshal(r, &o); err != nil {
			return AddResponse{}, err
		}
		if res, ok := o["result"]; ok {
			return res, nil
		} else {
			return AddResponse{}, errors.New("Unable to find result")
		}
	}

}
