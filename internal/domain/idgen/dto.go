package idgen

import "digit-oss/noc-services/internal/domain"

// IdGenerationRequest wraps a batch of IdRequest items.
type IdGenerationRequest struct {
	RequestInfo *domain.RequestInfo `json:"RequestInfo"`
	IdRequests  []IdRequest         `json:"idRequests"`
}

// IdGenerationResponse wraps a batch of IdResponse items.
type IdGenerationResponse struct {
	ResponseInfo *domain.ResponseInfo `json:"responseInfo"`
	IdResponses  []IdResponse         `json:"idResponses"`
}
