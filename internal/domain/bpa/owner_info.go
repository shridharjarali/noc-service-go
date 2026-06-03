package bpa

import (
	"encoding/json"

	"digit-oss/noc-services/internal/domain"
)

// OwnerInfo represents property owner information.
type OwnerInfo struct {
	TenantID              string               `json:"tenantId"`
	Name                  string               `json:"name"`
	OwnerID               string               `json:"ownerId"`
	MobileNumber          string               `json:"mobileNumber"`
	Gender                string               `json:"gender"`
	FatherOrHusbandName   string               `json:"fatherOrHusbandName"`
	CorrespondenceAddress string               `json:"correspondenceAddress"`
	IsPrimaryOwner        *bool                `json:"isPrimaryOwner"`
	Status                *bool                `json:"status"`
	OwnerShipPercentage   *float64             `json:"ownerShipPercentage"`
	OwnerType             string               `json:"ownerType"`
	InstitutionID         string               `json:"institutionId"`
	Documents             []Document           `json:"documents"`
	Relationship          Relationship         `json:"relationship"`
	AdditionalDetails     json.RawMessage      `json:"additionalDetails"`
	ID                    *int64               `json:"id"`
	UUID                  string               `json:"uuid"`
	UserName              string               `json:"userName"`
	Password              string               `json:"password"`
	Salutation            string               `json:"salutation"`
	EmailID               string               `json:"emailId"`
	AltContactNumber      string               `json:"altContactNumber"`
	Pan                   string               `json:"pan"`
	AadhaarNumber         string               `json:"aadhaarNumber"`
	PermanentAddress      string               `json:"permanentAddress"`
	PermanentCity         string               `json:"permanentCity"`
	PermanentPincode      string               `json:"permanentPinCode"`
	CorrespondenceCity    string               `json:"correspondenceCity"`
	CorrespondencePincode string               `json:"correspondencePinCode"`
	Active                *bool                `json:"active"`
	Dob                   *int64               `json:"dob"`
	PwdExpiryDate         *int64               `json:"pwdExpiryDate"`
	Locale                string               `json:"locale"`
	Type                  string               `json:"type"`
	Signature             string               `json:"signature"`
	AccountLocked         *bool                `json:"accountLocked"`
	Roles                 []domain.Role        `json:"roles"`
	BloodGroup            string               `json:"bloodGroup"`
	IdentificationMark    string               `json:"identificationMark"`
	Photo                 string               `json:"photo"`
	CreatedBy             string               `json:"createdBy"`
	CreatedDate           *int64               `json:"createdDate"`
	LastModifiedBy        string               `json:"lastModifiedBy"`
	LastModifiedDate      *int64               `json:"lastModifiedDate"`
	OtpReference          string               `json:"otpReference"`
	AuditDetails          *domain.AuditDetails `json:"auditDetails"`
}
