package domain

// Workflow captures workflow action data.
type Workflow struct {
	Action    string     `json:"action"`
	Assignes  []string   `json:"assignes"`
	Comment   string     `json:"comment"`
	Documents []Document `json:"documents"`
}
