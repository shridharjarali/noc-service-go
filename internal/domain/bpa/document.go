package bpa

import (
	"encoding/json"

	"digit-oss/noc-services/internal/domain"
)

// Document holds documents specific to BPA (includes auditDetails).
type Document struct {
	ID                string               `json:"id"`
	DocumentType      string               `json:"documentType"`
	FileStoreID       string               `json:"fileStoreId"`
	DocumentUID       string               `json:"documentUid"`
	AdditionalDetails json.RawMessage      `json:"additionalDetails"`
	AuditDetails      *domain.AuditDetails `json:"auditDetails"`
}
