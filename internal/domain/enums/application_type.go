package enums

// ApplicationType represents the type of a NOC application.
type ApplicationType string

const (
	ApplicationTypeProvisional ApplicationType = "PROVISIONAL"
	ApplicationTypeNew         ApplicationType = "NEW"
	ApplicationTypeRenew       ApplicationType = "RENEW"
)
