package idgen

// IdRequest is a single ID generation request item.
type IdRequest struct {
	IdName   string `json:"idName"`
	TenantID string `json:"tenantId"`
	Format   string `json:"format"`
}
