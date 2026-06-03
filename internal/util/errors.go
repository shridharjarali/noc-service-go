package util

import "fmt"

// CustomError mirrors org.egov.tracer.model.CustomException.
// Code is the error key; Message is the human-readable description.
type CustomError struct {
	Code    string
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewCustomError creates a single-code CustomError.
func NewCustomError(code, message string) *CustomError {
	return &CustomError{Code: code, Message: message}
}

// MultiError collects multiple validation errors.
type MultiError struct {
	Errors map[string]string
}

func (e *MultiError) Error() string {
	return fmt.Sprintf("validation errors: %v", e.Errors)
}

// NewMultiError creates a MultiError and returns it if non-empty, else nil.
func NewMultiError(errors map[string]string) error {
	if len(errors) == 0 {
		return nil
	}
	return &MultiError{Errors: errors}
}
