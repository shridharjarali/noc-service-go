package models

import (
	"encoding/json"

	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/domain/enums"
)

// ──────────────────────────────────────────────────────────────────────────────
// NocModel maps to the existing eg_noc table managed by egov-persister.
// GORM reads from this table; writes still go through Kafka.
// ──────────────────────────────────────────────────────────────────────────────

type NocModel struct {
	ID                string             `gorm:"column:id;primaryKey"`
	TenantID          string             `gorm:"column:tenantid"`
	ApplicationNo     string             `gorm:"column:applicationno"`
	NocNo             string             `gorm:"column:nocno"`
	NocType           string             `gorm:"column:noctype"`
	ApplicationType   string             `gorm:"column:applicationtype"`
	ApplicationStatus string             `gorm:"column:applicationstatus"`
	Status            string             `gorm:"column:status"`
	LandID            string             `gorm:"column:landid"`
	Source            string             `gorm:"column:source"`
	SourceRefID       string             `gorm:"column:sourcerefid"`
	AccountID         string             `gorm:"column:accountid"`
	AdditionalDetails string             `gorm:"column:additionaldetails;type:jsonb"`
	CreatedBy         string             `gorm:"column:createdby"`
	LastModifiedBy    string             `gorm:"column:lastmodifiedby"`
	CreatedTime       int64              `gorm:"column:createdtime"`
	LastModifiedTime  int64              `gorm:"column:lastmodifiedtime"`
	Documents         []NocDocumentModel `gorm:"foreignKey:NocID;references:ID"`
}

// TableName overrides GORM's default table name convention.
func (NocModel) TableName() string { return "eg_noc" }

// ──────────────────────────────────────────────────────────────────────────────
// NocDocumentModel maps to the existing eg_noc_document table.
// ──────────────────────────────────────────────────────────────────────────────

type NocDocumentModel struct {
	ID                string `gorm:"column:id;primaryKey"`
	NocID             string `gorm:"column:nocid"`
	DocumentType      string `gorm:"column:documenttype"`
	FileStoreID       string `gorm:"column:filestoreid"`
	DocumentUID       string `gorm:"column:documentuid"`
	AdditionalDetails string `gorm:"column:additionaldetails;type:jsonb"`
}

// TableName overrides GORM's default table name convention.
func (NocDocumentModel) TableName() string { return "eg_noc_document" }

// ──────────────────────────────────────────────────────────────────────────────
// Domain ↔ Model Converters
// ──────────────────────────────────────────────────────────────────────────────

// ToDomain converts a NocModel to the existing domain.Noc struct.
func (m *NocModel) ToDomain() domain.Noc {
	ct := m.CreatedTime
	lmt := m.LastModifiedTime

	var addlJSON json.RawMessage
	if m.AdditionalDetails != "" && m.AdditionalDetails != "{}" && m.AdditionalDetails != "null" {
		addlJSON = json.RawMessage(m.AdditionalDetails)
	}

	noc := domain.Noc{
		ID:                m.ID,
		TenantID:          m.TenantID,
		ApplicationNo:     m.ApplicationNo,
		NocNo:             m.NocNo,
		NocType:           m.NocType,
		ApplicationType:   enums.ApplicationType(m.ApplicationType),
		ApplicationStatus: m.ApplicationStatus,
		Status:            enums.Status(m.Status),
		LandID:            m.LandID,
		Source:            m.Source,
		SourceRefID:       m.SourceRefID,
		AccountID:         m.AccountID,
		AdditionalDetails: addlJSON,
		AuditDetails: &domain.AuditDetails{
			CreatedBy:        m.CreatedBy,
			LastModifiedBy:   m.LastModifiedBy,
			CreatedTime:      &ct,
			LastModifiedTime: &lmt,
		},
	}

	for i := range m.Documents {
		noc.Documents = append(noc.Documents, m.Documents[i].ToDomain())
	}

	return noc
}

// ToDomain converts a NocDocumentModel to the existing domain.Document struct.
func (d *NocDocumentModel) ToDomain() domain.Document {
	var docAddlJSON json.RawMessage
	if d.AdditionalDetails != "" && d.AdditionalDetails != "{}" && d.AdditionalDetails != "null" {
		docAddlJSON = json.RawMessage(d.AdditionalDetails)
	}

	return domain.Document{
		ID:                d.ID,
		DocumentType:      d.DocumentType,
		FileStoreID:       d.FileStoreID,
		DocumentUID:       d.DocumentUID,
		AdditionalDetails: docAddlJSON,
	}
}
