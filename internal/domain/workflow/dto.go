package workflow

import "digit-oss/noc-services/internal/domain"

// BusinessServiceResponse wraps a list of BusinessService.
type BusinessServiceResponse struct {
	ResponseInfo     *domain.ResponseInfo `json:"ResponseInfo"`
	BusinessServices []BusinessService    `json:"BusinessServices"`
}

// ProcessInstanceResponse wraps a list of ProcessInstance.
type ProcessInstanceResponse struct {
	ResponseInfo     *domain.ResponseInfo `json:"ResponseInfo"`
	ProcessInstances []ProcessInstance     `json:"ProcessInstances"`
}
