package yad

// YandexTokenbody is used in MakeYandexOauthRequest()
// to unmarshal yandex response body and get AccessToken
type YandexTokenbody struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

//TODO: upgrade type to common Yandex API error (DONE)
type YandexV5Error struct {
	ErrorDescription string `json:"error_description"`
	ErrorCode        string `json:"error_code"`
	ErrorString      string `json:"error_string"`
	RequsetID        string `json:"request_id"`
}

//"error_str":"Authorization error","error_code":53,"error_detail":""
type YandexV4Error struct {
	ErrorDescription string `json:"error_str"`
	ErrorDetail      string `json:"error_detail"`
	ErrorCode        int    `json:"error_code"`
}

//YandexDirectAPI V5 Error
// "error" : {
// "request_id": "8695244274068608439",
// "error_code": 54,
// "error_string": "No rights",
// "error_detail": "No rights to indicated client"
// }
// Python example
//"Ошибка API
// {$apiErr->error_code}: " +
//"{$apiErr->error_string} - " +
//"{$apiErr->error_detail} " +
//"(RequestId: {$apiErr->request_id})"

//YandexDirectAPI V4 Error
//, "API error:"
//, "Code: ".$result->{error_code}
//, "Describe: ".$result->{error_str}
//, "Detail:".($result->{error_detail} || '');

type GetSummaryStatRes struct {
	Data []StatisticDataType `json:"Data"`
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

//////

//////

type ParamType struct {
	CampaignIDS []string `json:"CampaignIDS"`
	StartDate   string   `json:"StartDate"`
	EndDate     string   `json:"EndDate"`
}
