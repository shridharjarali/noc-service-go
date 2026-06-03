package bpa

import "encoding/json"

// Institution represents an institutional property owner.
type Institution struct {
	ID                     string          `json:"id"`
	TenantID               string          `json:"tenantId"`
	Type                   string          `json:"type"`
	Designation            string          `json:"designation"`
	NameOfAuthorizedPerson string          `json:"nameOfAuthorizedPerson"`
	AdditionalDetails      json.RawMessage `json:"additionalDetails"`
}
