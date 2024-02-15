package httplib

import (
	"fmt"
	"net/http"
	internalErr "projects/fb-server/pkg/errors"
	"time"
)

// ApiError represents a structure for encoding API error responses in JSON format.
type ApiError struct {
	HttpStatus int `json:"http_status,omitempty"`
	ErrorCode  int `json:"code"`
	Message    any `json:"message"`
	Timestamp  any `json:"timestamp"`
}

// Error returns the string representation of the ApiError.
func (e ApiError) Error() string {
	return e.String()
}

// String returns a formatted string representation of the ApiError, including HTTP status code,
// internal error code, and the error message.
func (e ApiError) String() string {
	msg := fmt.Sprintf("%s", e.Message)

	if e.HttpStatus > 0 {
		msg = fmt.Sprintf("[HTTP_CODE:%d] %s", e.HttpStatus, msg)
	}

	if e.ErrorCode > 0 {
		msg = fmt.Sprintf("[INTERNAL_CODE:%d] %s", e.ErrorCode, msg)
	}

	return msg
}

// NewApiError creates and returns a new instance of ApiError with the specified error code and message.
// It automatically sets the timestamp to the current time in RFC1123 format.
func NewApiError(code int, msg string) *ApiError {
	return &ApiError{
		ErrorCode: code,
		Message:   msg,
		Timestamp: time.Now().Format(time.RFC1123),
	}
}

// NewApiErrFromInternalErr creates and returns a new instance of ApiError based on an internal error.
// It sets the HTTP status code to the specified value or defaults to http.StatusBadRequest if not provided.
func NewApiErrFromInternalErr(e *internalErr.InternalError, code ...int) *ApiError {
	var status int

	if len(code) > 0 {
		status = code[0]
	} else {
		status = http.StatusBadRequest
	}

	return &ApiError{
		HttpStatus: status,
		ErrorCode:  e.GetCode(),
		Message:    e.Message,
		Timestamp:  time.Now().Format(time.RFC1123),
	}
}

// CreatedObjectId represents a structure for encoding the ID and UUID of a created entity.
type CreatedObjectId struct {
	Id   any `json:"id,omitempty" yaml:"id,omitempty"`
	UUID any `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

// ListResult represents a structure for encoding API list of results in JSON format.
type ListResult struct {
	Results any   `json:"results,omitempty" yaml:"results,omitempty"`
	Count   int32 `json:"count" yaml:"count"`
}

// Response represents a structure for encoding API responses in JSON format.
type Response struct {
	Success   bool   `json:"success" yaml:"success"`
	Timestamp string `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
	Message   string `json:"message,omitempty" yaml:"message,omitempty"`
	Data      any    `json:"data,omitempty" yaml:"data,omitempty"`
	CreatedObjectId
	RequestStatus *ApiError `json:"request_status" yaml:"request_status,omitempty"`
}

// SuccessfulResult returns a Response instance with success flag set to true and default message.
func SuccessfulResult() Response {
	return Response{
		Success:   true,
		Timestamp: time.Now().Format(time.RFC1123),
		Message:   "The service is operating normally",
	}
}

// SuccessfulResultMap returns a map with success flag set to true and default message.
func SuccessfulResultMap() map[string]any {
	return map[string]any{
		"success":   true,
		"timestamp": time.Now().Format(time.RFC1123),
		"message":   "The service is operating normally",
	}
}
