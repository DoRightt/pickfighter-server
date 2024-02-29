package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ApiServiceMock struct {
	ServiceName string
}

func (m *ApiServiceMock) Init(ctx context.Context) error           { return nil }
func (m *ApiServiceMock) ApplyRoutes()                             {}
func (m *ApiServiceMock) Shutdown(ctx context.Context, sig string) {}

func TestHealthCheck(t *testing.T) {
	apiHandler := &ApiHandler{
		ServiceName: "TestApp",
		Services:    make(map[string]ApiService),
	}

	apiHandler.Services["TestA"] = &ApiServiceMock{ServiceName: "TestA"}
	apiHandler.Services["TestB"] = &ApiServiceMock{ServiceName: "TestB"}

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	apiHandler.HealthCheck(w, req)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.ElementsMatch(t, []interface{}{"TestA", "TestB"}, response["modules"])
}
