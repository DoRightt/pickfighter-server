package httplib

import (
	"net/http"
	internalErr "projects/fb-server/pkg/errors"
	"testing"
	"time"
)

func TestApiError(t *testing.T) {
	testCases := []struct {
		Name           string
		ApiErr         ApiError
		ExpectedString string
	}{
		{
			Name: "Full ApiError",
			ApiErr: ApiError{
				HttpStatus: 404,
				ErrorCode:  1001,
				Message:    "Not Found",
				Timestamp:  time.Now(),
			},
			ExpectedString: "[INTERNAL_CODE:1001] [HTTP_CODE:404] Not Found",
		},
		{
			Name: "ApiError without HttpStatus",
			ApiErr: ApiError{
				ErrorCode: 1002,
				Message:   "Error without HTTP status",
				Timestamp: time.Now(),
			},
			ExpectedString: "[INTERNAL_CODE:1002] Error without HTTP status",
		},
		{
			Name: "ApiError without ErrorCode",
			ApiErr: ApiError{
				HttpStatus: 500,
				Message:    "Error without internal code",
				Timestamp:  time.Now(),
			},
			ExpectedString: "[HTTP_CODE:500] Error without internal code",
		},
		{
			Name: "ApiError without HttpStatus and ErrorCode",
			ApiErr: ApiError{
				Message:   "Error without HTTP status and internal code",
				Timestamp: time.Now(),
			},
			ExpectedString: "Error without HTTP status and internal code",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			resultString := tc.ApiErr.String()
			if resultString != tc.ExpectedString {
				t.Errorf("Expected %v, but got %v", tc.ExpectedString, resultString)
			}

			resultError := tc.ApiErr.Error()
			if resultError != tc.ExpectedString {
				t.Errorf("Expected %v, but got %v", tc.ExpectedString, resultError)
			}
		})
	}
}

func TestNewApiError(t *testing.T) {
	code := 404
	msg := "Not Found"

	apiError := NewApiError(code, msg)

	if apiError == nil {
		t.Errorf("NewApiError(%d, %s) returned nil", code, msg)
	}

	if apiError.ErrorCode != code {
		t.Errorf("NewApiError(%d, %s) ErrorCode = %d; expected %d", code, msg, apiError.ErrorCode, code)
	}
	if apiError.Message != msg {
		t.Errorf("NewApiError(%d, %s) Message = %s; expected %s", code, msg, apiError.Message, msg)
	}

	if apiError.Timestamp == "" {
		t.Error("NewApiError() Timestamp is an empty string")
	}

	_, err := time.Parse(time.RFC1123, apiError.Timestamp.(string))
	if err != nil {
		t.Errorf("NewApiError() Timestamp has an invalid format: %v", err)
	}
}

func TestNewApiErrorWithCode(t *testing.T) {
	internalErr := &internalErr.InternalError{
		Code:    123,
		Message: "Test internal error",
	}

	apiError := NewApiErrFromInternalErr(internalErr)
	apiErrorWithCode := NewApiErrFromInternalErr(internalErr, http.StatusOK)

	if apiError == nil {
		t.Error("NewApiErrFromInternalErr() returned nil")
	}

	if apiError.HttpStatus != http.StatusBadRequest {
		t.Errorf("NewApiErrFromInternalErr() HttpStatus = %d; expected %d", apiError.HttpStatus, http.StatusBadRequest)
	}
	if apiErrorWithCode.HttpStatus != http.StatusOK {
		t.Errorf("NewApiErrFromInternalErr() HttpStatus = %d; expected %d", apiError.HttpStatus, http.StatusBadRequest)
	}
	if apiError.ErrorCode != internalErr.Code {
		t.Errorf("NewApiErrFromInternalErr() ErrorCode = %d; expected %d", apiError.ErrorCode, internalErr.Code)
	}
	if apiError.Message != internalErr.Message {
		t.Errorf("NewApiErrFromInternalErr() Message = %s; expected %s", apiError.Message, internalErr.Message)
	}

	if apiError.Timestamp == "" {
		t.Error("NewApiErrFromInternalErr() Timestamp is an empty string")
	}

	_, err := time.Parse(time.RFC1123, apiError.Timestamp.(string))
	if err != nil {
		t.Errorf("NewApiErrFromInternalErr() Timestamp has an invalid format: %v", err)
	}
}

func TestSuccessfulResult(t *testing.T) {
	response := SuccessfulResult()

	if response.Success != true {
		t.Fatalf("expected %v but got %v", true, response.Success)
	}

	_, err := time.Parse(time.RFC1123, response.Timestamp)
	if err != nil {
		t.Fatalf("expected time format to be RFC1123 but got error: %v", err)
	}

	expectedMessage := "The service is operating normally"
	if response.Message != expectedMessage {
		t.Fatalf("expected %v but got %v", expectedMessage, response.Message)
	}
}

func TestSuccessfulResultMap(t *testing.T) {
	response := SuccessfulResultMap()

	if response["success"] != true {
		t.Fatalf("expected %v but got %v", true, response["success"])
	}

	_, err := time.Parse(time.RFC1123, response["timestamp"].(string))
	if err != nil {
		t.Fatalf("expected time format to be RFC1123 but got error: %v", err)
	}

	expectedMessage := "The service is operating normally"
	if response["message"] != expectedMessage {
		t.Fatalf("expected %v but got %v", expectedMessage, response["message"])
	}
}
