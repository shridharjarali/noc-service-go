package bpa

import "digit-oss/noc-services/internal/domain"

// BPARequest is the inbound request contract for BPA operations.
type BPARequest struct {
	RequestInfo *domain.RequestInfo `json:"RequestInfo"`
	BPA         *BPA                `json:"BPA"`
}

// BPAResponse is the outbound response contract for BPA operations.
type BPAResponse struct {
	ResponseInfo *domain.ResponseInfo `json:"ResponseInfo"`
	BPA          []BPA                `json:"BPA"`
	Count        int                  `json:"Count"`
}

// BPASearchCriteria holds query parameters for searching BPA applications.
type BPASearchCriteria struct {
	TenantID        string   `json:"tenantId"`
	IDs             []string `json:"ids"`
	Status          string   `json:"status"`
	EdcrNumber      string   `json:"edcrNumber"`
	ApplicationNo   string   `json:"applicationNo"`
	ApprovalNo      string   `json:"approvalNo"`
	MobileNumber    string   `json:"mobileNumber"`
	LandID          []string `json:"-"`
	Offset          *int     `json:"offset"`
	Limit           *int     `json:"limit"`
	ApprovalDate    *int64   `json:"approvalDate"`
	FromDate        *int64   `json:"fromDate"`
	ToDate          *int64   `json:"toDate"`
	OwnerIDs        []string `json:"-"`
	BusinessService []string `json:"-"`
	CreatedBy       []string `json:"-"`
	Locality        string   `json:"locality"`
	ApplicationType string   `json:"applicationType"`
	ServiceType     string   `json:"serviceType"`
	PermitNumber    string   `json:"permitNumber"`
}
