package domain

// AuditDetails captures audit-related fields.
type AuditDetails struct {
	CreatedBy        string `json:"createdBy"`
	LastModifiedBy   string `json:"lastModifiedBy"`
	CreatedTime      *int64 `json:"createdTime"`
	LastModifiedTime *int64 `json:"lastModifiedTime"`
}
