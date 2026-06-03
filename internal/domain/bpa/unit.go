package bpa

import (
	"encoding/json"

	"digit-oss/noc-services/internal/domain"
)

// Unit represents a building unit.
type Unit struct {
	ID                string               `json:"id"`
	TenantID          string               `json:"tenantId"`
	FloorNo           string               `json:"floorNo"`
	UnitType          string               `json:"unitType"`
	UsageCategory     string               `json:"usageCategory"`
	OccupancyType     string               `json:"occupancyType"`
	OccupancyDate     *int64               `json:"occupancyDate"`
	AdditionalDetails json.RawMessage      `json:"additionalDetails"`
	AuditDetails      *domain.AuditDetails `json:"auditDetails"`
}
