package httplib

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeError(t *testing.T) {
	internalCode := 500
	err := errors.New("Internal server error")

	apiErr := makeError(internalCode, err)

	assert.NotNil(t, apiErr, "makeError should return a non-nil ApiError")
	assert.Equal(t, internalCode, apiErr.ErrorCode, "ErrorCode should match the provided internalCode")
	assert.Equal(t, err.Error(), apiErr.Message, "Message should match the error message")
}

func TestErrorResponseJSON(t *testing.T) {
	internalCode := 500
	err := errors.New("Internal server error")

	w := httptest.NewRecorder()

	ErrorResponseJSON(w, http.StatusInternalServerError, internalCode, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Status code should be http.StatusInternalServerError")

	var apiErr ApiError
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &apiErr), "Error decoding JSON response")
	assert.Equal(t, internalCode, apiErr.ErrorCode, "ErrorCode should match the provided internalCode")
	assert.Equal(t, err.Error(), apiErr.Message, "Message should match the error message")
}

func TestResponseJSON(t *testing.T) {
	testData := map[string]interface{}{
		"key": "value",
	}

	w := httptest.NewRecorder()

	ResponseJSON(w, testData)

	assert.Equal(t, http.StatusOK, w.Code, "Status code should be http.StatusOK")

	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response), "Error decoding JSON response")
	assert.Equal(t, testData, response, "Response data should match the provided testData")
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name       string
		httpCode   int
		v          interface{}
		wantHeader http.Header
		wantBody   string
	}{
		{
			name:     "Write Success JSON",
			httpCode: 200,
			v:        map[string]string{"message": "Success"},
			wantHeader: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			wantBody: `{"message":"Success"}`,
		},
		{
			name:     "Write Error JSON",
			httpCode: 500,
			v:        map[string]string{"message": "Error"},
			wantHeader: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			wantBody: `{"message":"Error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			writeJSON(recorder, tc.httpCode, tc.v)

			res := recorder.Result()

			if !reflect.DeepEqual(res.Header, tc.wantHeader) {
				t.Errorf("writeJSON header = %v, want %v", res.Header, tc.wantHeader)
			}

			if body, _ := json.Marshal(tc.v); string(body) != tc.wantBody {
				t.Errorf("writeJSON body = %v, want %v", string(body), tc.wantBody)
			}
		})
	}
}
