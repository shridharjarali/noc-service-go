package workflow

import "digit-oss/noc-services/internal/domain"

// BusinessService represents a workflow business service definition.
type BusinessService struct {
	TenantID           string               `json:"tenantId"`
	UUID               string               `json:"uuid"`
	BusinessService    string               `json:"businessService"`
	Business           string               `json:"business"`
	GetURI             string               `json:"getUri"`
	PostURI            string               `json:"postUri"`
	BusinessServiceSla *int64               `json:"businessServiceSla"`
	States             []State              `json:"states"`
	AuditDetails       *domain.AuditDetails `json:"auditDetails"`
}
