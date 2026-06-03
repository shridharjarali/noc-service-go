package util

import "time"

import "digit-oss/noc-services/internal/domain"

// GetAuditDetails builds AuditDetails for create or update flows.
// Mirrors NOCUtil.getAuditDetails in Java.
func GetAuditDetails(userUUID string, isCreate bool) *domain.AuditDetails {
	now := TimeNowMillis()
	if isCreate {
		return &domain.AuditDetails{
			CreatedBy:        userUUID,
			LastModifiedBy:   userUUID,
			CreatedTime:      &now,
			LastModifiedTime: &now,
		}
	}
	return &domain.AuditDetails{
		LastModifiedBy:   userUUID,
		LastModifiedTime: &now,
	}
}

// TimeNowMillis returns current time in epoch milliseconds.
func TimeNowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
