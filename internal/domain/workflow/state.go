package workflow

import "digit-oss/noc-services/internal/domain"

// State represents a state in a business-service workflow.
type State struct {
	UUID              string               `json:"uuid"`
	TenantID          string               `json:"tenantId"`
	BusinessServiceID string               `json:"businessServiceId"`
	Sla               *int64               `json:"sla"`
	State             string               `json:"state"`
	ApplicationStatus string               `json:"applicationStatus"`
	DocUploadRequired *bool                `json:"docUploadRequired"`
	IsStartState      *bool                `json:"isStartState"`
	IsTerminateState  *bool                `json:"isTerminateState"`
	IsStateUpdatable  *bool                `json:"isStateUpdatable"`
	Actions           []Action             `json:"actions"`
	AuditDetails      *domain.AuditDetails `json:"auditDetails"`
}
