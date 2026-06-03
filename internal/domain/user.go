package domain

// UserSearchResponse represents a user record returned from user search.
type UserSearchResponse struct {
	ID                    *int64  `json:"id"`
	UUID                  string  `json:"uuid"`
	UserName              string  `json:"userName"`
	Password              string  `json:"password"`
	Salutation            string  `json:"salutation"`
	Name                  string  `json:"name"`
	Gender                string  `json:"gender"`
	MobileNumber          string  `json:"mobileNumber"`
	EmailID               string  `json:"emailId"`
	AltContactNumber      string  `json:"altContactNumber"`
	Pan                   string  `json:"pan"`
	AadhaarNumber         string  `json:"aadhaarNumber"`
	PermanentAddress      string  `json:"permanentAddress"`
	PermanentCity         string  `json:"permanentCity"`
	PermanentPincode      string  `json:"permanentPinCode"`
	CorrespondenceCity    string  `json:"correspondenceCity"`
	CorrespondencePincode string  `json:"correspondencePinCode"`
	CorrespondenceAddress string  `json:"correspondenceAddress"`
	Active                *bool   `json:"active"`
	Dob                   *int64  `json:"dob"`
	PwdExpiryDate         *int64  `json:"pwdExpiryDate"`
	Locale                string  `json:"locale"`
	Type                  string  `json:"type"`
	Signature             string  `json:"signature"`
	AccountLocked         *bool   `json:"accountLocked"`
	FatherOrHusbandName   string  `json:"fatherOrHusbandName"`
	BloodGroup            string  `json:"bloodGroup"`
	IdentificationMark    string  `json:"identificationMark"`
	Photo                 string  `json:"photo"`
	CreatedBy             string  `json:"createdBy"`
	CreatedDate           *int64  `json:"createdDate"`
	LastModifiedBy        string  `json:"lastModifiedBy"`
	LastModifiedDate      *int64  `json:"lastModifiedDate"`
	OtpReference          string  `json:"otpReference"`
	TenantID              string  `json:"tenantId"`
}

// User mirrors org.egov.common.contract.request.User.
type User struct {
	ID           int64  `json:"id"`
	UUID         string `json:"uuid"`
	UserName     string `json:"userName"`
	Name         string `json:"name"`
	MobileNumber string `json:"mobileNumber"`
	EmailID      string `json:"emailId"`
	TenantID     string `json:"tenantId"`
	Roles        []Role `json:"roles"`
	Type         string `json:"type"`
}

// Role represents a user role (mirrors org.egov.common.contract.request.Role).
type Role struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	TenantID string `json:"tenantId"`
}
