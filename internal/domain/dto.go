package domain

// ─── NOC Request / Response / Search ─────────────────────────────────────────

// NocRequest binds the metadata contract and the main NOC application contract.
type NocRequest struct {
	RequestInfo *RequestInfo `json:"RequestInfo"`
	Noc         *Noc         `json:"Noc"`
}

// NocResponse contains the response metadata and a list of NOC applications.
type NocResponse struct {
	ResponseInfo *ResponseInfo `json:"ResponseInfo"`
	Noc          []Noc         `json:"Noc"`
	Count        *int          `json:"count"`
}

// NocSearchCriteria holds the query parameters for searching NOC applications.
type NocSearchCriteria struct {
	TenantID      string   `json:"tenantId"`
	IDs           []string `json:"ids"`
	ApplicationNo string   `json:"applicationNo"`
	MobileNumber  string   `json:"mobileNumber"`
	NocNo         string   `json:"nocNo"`
	Source        string   `json:"source"`
	NocType       string   `json:"nocType"`
	SourceRefID   string   `json:"sourceRefId"`
	Offset        *int     `json:"offset"`
	Limit         *int     `json:"limit"`
	OwnerIDs      []string `json:"-"`
	AccountID     []string `json:"accountId"`
	Status        []string `json:"status"`
}

// RequestInfoWrapper wraps RequestInfo for endpoints that only need metadata.
type RequestInfoWrapper struct {
	RequestInfo *RequestInfo `json:"RequestInfo"`
}

// ─── User DTOs ───────────────────────────────────────────────────────────────

// UserSearchRequest is the payload sent to the user-search API.
type UserSearchRequest struct {
	RequestInfo   *RequestInfo `json:"RequestInfo"`
	UUID          []string     `json:"uuid"`
	ID            []string     `json:"id"`
	UserName      string       `json:"userName"`
	Name          string       `json:"name"`
	MobileNumber  string       `json:"mobileNumber"`
	AadhaarNumber string       `json:"aadhaarNumber"`
	Pan           string       `json:"pan"`
	EmailID       string       `json:"emailId"`
	FuzzyLogic    bool         `json:"fuzzyLogic"`
	Active        *bool        `json:"active"`
	TenantID      string       `json:"tenantId"`
	PageSize      int          `json:"pageSize"`
	PageNumber    int          `json:"pageNumber"`
	Sort          []string     `json:"sort"`
	UserType      string       `json:"userType"`
	RoleCodes     []string     `json:"roleCodes"`
}

// UserResponse wraps the response from user-search.
type UserResponse struct {
	ResponseInfo *ResponseInfo        `json:"responseInfo"`
	User         []UserSearchResponse `json:"user"`
}
