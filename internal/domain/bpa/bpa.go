package bpa

import (
	"encoding/json"

	"digit-oss/noc-services/internal/domain"
)

// BPA is the Building Plan Approval application object.
type BPA struct {
	ID                string              `json:"id"`
	ApplicationNo     string              `json:"applicationNo"`
	ApprovalNo        string              `json:"approvalNo"`
	AccountID         string              `json:"accountId"`
	EdcrNumber        string              `json:"edcrNumber"`
	RiskType          string              `json:"riskType"`
	BusinessService   string              `json:"businessService"`
	LandID            string              `json:"landId"`
	TenantID          string              `json:"tenantId"`
	ApprovalDate      *int64              `json:"approvalDate"`
	ApplicationDate   *int64              `json:"applicationDate"`
	Status            string              `json:"status"`
	Documents         []Document          `json:"documents"`
	LandInfo          *LandInfo           `json:"landInfo"`
	Workflow          *domain.Workflow     `json:"workflow"`
	AuditDetails      *domain.AuditDetails `json:"auditDetails"`
	AdditionalDetails json.RawMessage     `json:"additionalDetails"`
}
