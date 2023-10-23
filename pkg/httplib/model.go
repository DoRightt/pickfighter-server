package httplib

import "time"

type Response struct {
	Success   bool   `json:"success" yaml:"success"`
	Timestamp string `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
	Message   string `json:"message,omitempty" yaml:"message,omitempty"`
	Data      any    `json:"data,omitempty" yaml:"data,omitempty"`
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
