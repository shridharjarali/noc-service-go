package service

import (
	"digit-oss/noc-services/internal/config"
	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/repository/postgres"
	"digit-oss/noc-services/internal/util"
)

// UserService translates UserService.java.
// Calls the user-service /users/_search endpoint.
type UserService struct {
	Cfg        *config.Config
	SvcRequest *postgres.ServiceRequestRepository
}

// GetUser searches users by the ownerIds in NocSearchCriteria.
func (u *UserService) GetUser(criteria domain.NocSearchCriteria, requestInfo *domain.RequestInfo) (*domain.UserResponse, error) {
	active := true
	tenantID := criteria.TenantID
	// Use root tenant
	if parts := splitDot(tenantID); len(parts) > 0 {
		tenantID = parts[0]
	}

	userReq := domain.UserSearchRequest{
		RequestInfo: requestInfo,
		TenantID:    tenantID,
		Active:      &active,
	}
	if len(criteria.OwnerIDs) > 0 {
		userReq.UUID = criteria.OwnerIDs
	}

	url := u.Cfg.UserHost + u.Cfg.UserSearchEndpoint

	var resp domain.UserResponse
	if err := u.SvcRequest.FetchResult(url, userReq, &resp); err != nil {
		return nil, util.NewCustomError("USER_SEARCH_ERROR", "Failed to search users: "+err.Error())
	}
	return &resp, nil
}

func splitDot(s string) []string {
	result := make([]string, 0)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}
