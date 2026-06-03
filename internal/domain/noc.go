package domain

import (
	"encoding/json"

	"digit-oss/noc-services/internal/domain/enums"
)

// Noc is the NOC application object capturing noc-related information,
// land ID, and related documents.
type Noc struct {
	ID                string               `json:"id"`
	TenantID          string               `json:"tenantId"`
	ApplicationNo     string               `json:"applicationNo"`
	NocNo             string               `json:"nocNo"`
	ApplicationType   enums.ApplicationType `json:"applicationType"`
	NocType           string               `json:"nocType"`
	AccountID         string               `json:"accountId"`
	Source            string               `json:"source"`
	SourceRefID       string               `json:"sourceRefId"`
	LandID            string               `json:"landId"`
	Status            enums.Status         `json:"status"`
	ApplicationStatus string               `json:"applicationStatus"`
	Documents         []Document           `json:"documents"`
	Workflow          *Workflow            `json:"workflow"`
	AuditDetails      *AuditDetails        `json:"auditDetails"`
	AdditionalDetails json.RawMessage      `json:"additionalDetails"`
}
