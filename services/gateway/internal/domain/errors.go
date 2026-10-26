package domain

import "fmt"

// AppException is the base interface for domain-specific errors.
type AppException struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func (e *AppException) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppException) Unwrap() error {
	return e.Err
}

// NewAppException creates a generic domain error.
func NewAppException(code string, msg string, status int, err error) *AppException {
	return &AppException{
		Code:       code,
		Message:    msg,
		StatusCode: status,
		Err:        err,
	}
}
