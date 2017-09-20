package clients

import "gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI/support"
import "gogs.itcloud.pro/SAS-project/sas/yandexDirectAPI/gc"

type GetRequest struct {
	support.ApiObject

	FieldNames []ClientFieldEnum `json:"FieldNames"`
}

type GetResponse struct {
	support.ApiObject

	Clients []gc.ClientGetItem `json:"Clients,omitempty"`
}
