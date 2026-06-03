package domain

// SMSRequest represents an SMS notification request.
type SMSRequest struct {
	MobileNumber string `json:"mobileNumber"`
	Message      string `json:"message"`
}
