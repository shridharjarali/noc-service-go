package domain

// RequestInfo mirrors org.egov.common.contract.request.RequestInfo.
type RequestInfo struct {
	APIId         string `json:"apiId"`
	Ver           string `json:"ver"`
	Ts            *int64 `json:"ts"`
	Action        string `json:"action"`
	Did           string `json:"did"`
	Key           string `json:"key"`
	MsgId         string `json:"msgId"`
	AuthToken     string `json:"authToken"`
	CorrelationId string `json:"correlationId"`
	UserInfo      *User  `json:"userInfo"`
}

// ResponseInfo mirrors org.egov.common.contract.response.ResponseInfo.
type ResponseInfo struct {
	APIId  string `json:"apiId"`
	Ver    string `json:"ver"`
	Ts     *int64 `json:"ts"`
	ResMsg string `json:"resMsgId"`
	MsgId  string `json:"msgId"`
	Status string `json:"status"`
}
