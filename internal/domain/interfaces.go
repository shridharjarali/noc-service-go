package domain

import "context"

// NOCRepository defines persistence operations for NOC applications.
// Save and Update push to Kafka (persister pattern), matching the Java version.
// Search queries PostgreSQL directly.
type NOCRepository interface {
	// Save pushes the NocRequest to the Kafka save topic.
	Save(ctx context.Context, req *NocRequest) error

	// Update pushes the NocRequest to the Kafka update topic.
	// If isStateUpdatable is true, it pushes to the update topic;
	// otherwise it pushes to the update-workflow topic.
	Update(ctx context.Context, req *NocRequest, isStateUpdatable bool) error

	// Search finds NOC applications matching the given criteria.
	Search(ctx context.Context, criteria NocSearchCriteria) ([]Noc, int, error)
}

// NOCService defines business-logic operations for NOC applications.
type NOCService interface {
	// Create validates and persists a new NOC application.
	Create(ctx context.Context, req *NocRequest) (*Noc, error)

	// Update validates and applies changes to an existing NOC application.
	Update(ctx context.Context, req *NocRequest) (*Noc, error)

	// Search retrieves NOC applications matching the criteria.
	Search(ctx context.Context, criteria NocSearchCriteria, requestInfo *RequestInfo) ([]Noc, int, error)
}
