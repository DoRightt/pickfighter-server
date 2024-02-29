package services_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"projects/fb-server/internal/services"
	mock_services "projects/fb-server/internal/services/mocks"
	mock_logger "projects/fb-server/pkg/logger/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)

	service := mock_services.NewMockApiService(ctrl)
	logger := mock_logger.NewMockFbLogger(ctrl)
	app := &services.ApiHandler{
		Logger:   logger,
		Services: map[string]services.ApiService{"mockService": service},
	}

	logger.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()

	var panicErr error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = fmt.Errorf("Panic: %v", r)
			}
		}()

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err)

		req.Header.Set("Origin", "http://example.com")

		app.ServeHTTP(w, req)

		expectedOrigin := "http://example.com"
		actualOrigin := w.Header().Get("Access-Control-Allow-Origin")

		assert.Equal(t, actualOrigin, expectedOrigin, fmt.Sprintf("Expected Access-Control-Allow-Origin: %s, got: %s", expectedOrigin, actualOrigin))

		expectedStatusCode := http.StatusOK
		actualStatusCode := w.Code

		assert.Equal(t, actualStatusCode, expectedStatusCode, fmt.Sprintf("Expected status code: %d, got: %d", expectedStatusCode, actualStatusCode))
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = fmt.Errorf("Panic: %v", r)
			}
		}()

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodOptions, "/test", nil)
		require.NoError(t, err)

		req.Header.Set("Origin", "http://example.com")

		app.ServeHTTP(w, req)

		expectedStatusCode := http.StatusOK
		actualStatusCode := w.Code

		assert.Equal(t, actualStatusCode, expectedStatusCode, fmt.Sprintf("Expected status code: %d, got: %d", expectedStatusCode, actualStatusCode))
	}()

	assert.NoError(t, panicErr)

}

func TestRunHTTPServer(t *testing.T) {
	// TODO
}
