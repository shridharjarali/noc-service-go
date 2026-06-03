package bpa

import (
	"encoding/json"

	"digit-oss/noc-services/internal/domain"
)

// LandInfo represents land information associated with a BPA.
type LandInfo struct {
	ID                string               `json:"id"`
	LandUID           string               `json:"landUId"`
	LandUniqueRegNo   string               `json:"landUniqueRegNo"`
	TenantID          string               `json:"tenantId"`
	Status            BPAStatus            `json:"status"`
	Address           *Address             `json:"address"`
	OwnershipCategory string               `json:"ownershipCategory"`
	Owners            []OwnerInfo          `json:"owners"`
	Institution       *Institution         `json:"institution"`
	Source            Source               `json:"source"`
	Channel           Channel              `json:"channel"`
	Documents         []Document           `json:"documents"`
	Unit              []Unit               `json:"unit"`
	AdditionalDetails json.RawMessage      `json:"additionalDetails"`
	AuditDetails      *domain.AuditDetails `json:"auditDetails"`
}
