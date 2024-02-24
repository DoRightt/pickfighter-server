package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	apiHandler := &ApiHandler{
		ServiceName: "TestApp",
	}

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	apiHandler.HealthCheck(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
