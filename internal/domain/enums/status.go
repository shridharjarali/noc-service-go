package enums

// Status represents the active/inactive status of a NOC record.
type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)
