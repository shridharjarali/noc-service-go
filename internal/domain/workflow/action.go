package workflow

import "digit-oss/noc-services/internal/domain"

// Action represents an action within a workflow state.
type Action struct {
	UUID         string               `json:"uuid"`
	TenantID     string               `json:"tenantId"`
	CurrentState string               `json:"currentState"`
	Action       string               `json:"action"`
	NextState    string               `json:"nextState"`
	Roles        []string             `json:"roles"`
	AuditDetails *domain.AuditDetails `json:"auditDetails"`
}
