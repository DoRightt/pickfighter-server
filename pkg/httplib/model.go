package httplib

import (
	"fmt"
	"net/http"
	internalErr "projects/fb-server/errors"
	"time"
)

type ApiError struct {
	HttpStatus int `json:"http_status,omitempty"`
	ErrorCode  int `json:"code"`
	Message    any `json:"message"`
	Timestamp  any `json:"timestamp"`
}

func (e ApiError) Error() string {
	return e.String()
}

func (e ApiError) String() string {
	msg := fmt.Sprintf("%s", e.Message)

	if e.HttpStatus > 0 {
		msg = fmt.Sprintf("[HTTP_CODE:%d] %s", e.HttpStatus, msg)
	}

	if e.ErrorCode > 0 {
		msg = fmt.Sprintf("[INTERNAL_CODE:%d] %s", e.ErrorCode, msg)
	}

	return fmt.Sprintf("%s", msg)
}

func NewApiError(code int, msg string) *ApiError {
	return &ApiError{
		ErrorCode: code,
		Message:   msg,
		Timestamp: time.Now().Format(time.RFC1123),
	}
}

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

type CreatedObjectId struct {
	Id   any `json:"id,omitempty" yaml:"id,omitempty"`
	UUID any `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type ListResult struct {
	Results any   `json:"results,omitempty" yaml:"results,omitempty"`
	Count   int32 `json:"count" yaml:"count"`
}

type Response struct {
	Success   bool   `json:"success" yaml:"success"`
	Timestamp string `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
	Message   string `json:"message,omitempty" yaml:"message,omitempty"`
	Data      any    `json:"data,omitempty" yaml:"data,omitempty"`
	CreatedObjectId
}

func SuccessfulResult() Response {
	return Response{
		Success:   true,
		Timestamp: time.Now().Format(time.RFC1123),
		Message:   "The service is operating normally",
	}
}

func SuccessfulResultMap() map[string]any {
	return map[string]any{
		"success":   true,
		"timestamp": time.Now().Format(time.RFC1123),
		"message":   "The service is operating normally",
	}
}
