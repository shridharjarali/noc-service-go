package workflow

import (
	"encoding/json"

	"digit-oss/noc-services/internal/domain"
)

// ProcessInstance represents a running workflow process instance.
type ProcessInstance struct {
	ID                  string               `json:"id"`
	TenantID            string               `json:"tenantId"`
	BusinessService     string               `json:"businessService"`
	BusinessID          string               `json:"businessId"`
	Action              string               `json:"action"`
	ModuleName          string               `json:"moduleName"`
	State               *State               `json:"state"`
	Comment             string               `json:"comment"`
	Documents           []domain.Document    `json:"documents"`
	Assigner            *domain.User         `json:"assigner"`
	Assignee            *domain.User         `json:"assignee"`
	NextActions         []Action             `json:"nextActions"`
	StateSla            *int64               `json:"stateSla"`
	BusinesssServiceSla *int64               `json:"businesssServiceSla"`
	PreviousStatus      string               `json:"previousStatus"`
	Entity              json.RawMessage      `json:"entity"`
	AuditDetails        *domain.AuditDetails `json:"auditDetails"`
}
