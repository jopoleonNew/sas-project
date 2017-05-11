package yad

type RequestV4 struct {
	Method  string   `json:"method"`
	Token   string   `json:"token"`
	ParamV4 *ParamV4 `json:"param"`
}

type ParamV4 struct {
	CampaignIDS []int  `json:"CampaignIDS"`
	StartDate   string `json:"StartDate"`
	EndDate     string `json:"EndDate"`
}

type RequestV5 struct {
	Method   string    `json:"method"`
	ParamsV5 *ParamsV5 `json:"params"`
}

type ParamsV5 struct {
	SelectionCriteria struct{} `json:"SelectionCriteria"`
	FieldNames        []string `json:"FieldNames"`
}
