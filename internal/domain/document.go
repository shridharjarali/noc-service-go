package domain

import "encoding/json"

// Document holds documents attached during a transaction.
type Document struct {
	ID                string          `json:"id"`
	DocumentType      string          `json:"documentType"`
	FileStoreID       string          `json:"fileStoreId"`
	DocumentUID       string          `json:"documentUid"`
	AdditionalDetails json.RawMessage `json:"additionalDetails"`
}
