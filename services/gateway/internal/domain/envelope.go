package domain

import "time"

// ResponseMeta holds standard API metadata.
type ResponseMeta struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

// APIResponse standardizes all gateway JSON outputs.
type APIResponse struct {
	Success bool              `json:"success"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Data    any               `json:"data,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
	Meta    ResponseMeta      `json:"meta"`
}

// NewSuccessResponse creates a successful response envelope.
func NewSuccessResponse(code string, msg string, data any, reqID string) APIResponse {
	return APIResponse{
		Success: true,
		Code:    code,
		Message: msg,
		Data:    data,
		Meta: ResponseMeta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			RequestID: reqID,
		},
	}
}

// NewErrorResponse creates a structured failure response envelope.
func NewErrorResponse(code string, msg string, errs map[string]string, reqID string) APIResponse {
	return APIResponse{
		Success: false,
		Code:    code,
		Message: msg,
		Errors:  errs,
		Meta: ResponseMeta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			RequestID: reqID,
		},
	}
}
