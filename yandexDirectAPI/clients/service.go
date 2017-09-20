package clients

import (
	"encoding/json"
	"errors"

	"gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI"
)

type Service struct {
	Client *yad.Client
}

func New(c *yad.Client) Service {
	return Service{c}
}

func (s Service) Get(p GetRequest) (GetResponse, error) {
	if r, err := s.Client.Do("clients", "get", p); err != nil {
		return GetResponse{}, err
	} else {
		o := map[string]GetResponse{}
		if err := json.Unmarshal(r, &o); err != nil {
			return GetResponse{}, err
		}
		if res, ok := o["result"]; ok {
			return res, nil
		} else {
			return GetResponse{}, errors.New("Unable to find result")
		}
	}

}
