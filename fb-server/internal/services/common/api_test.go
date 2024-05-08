package common

import (
	"context"
	"fmt"
	"fightbettr.com/fb-server/internal/services"
	mock_services "fightbettr.com/fb-server/internal/services/mocks"
	mock_logger "fightbettr.com/fb-server/pkg/logger/mocks"
	mock_pgxs "fightbettr.com/fb-server/pkg/pgxs/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockRepo := mock_pgxs.NewMockFbRepo(ctl)
	mockLogger := mock_logger.NewMockFbLogger(ctl)
	serviceName := "TestHandler"

	h := &services.ApiHandler{
		ServiceName: serviceName,
		Repo:        mockRepo,
		Logger:      mockLogger,
	}

	service := New(h)

	assert.NotNil(t, service, "Service must be initialized")
}

func TestInit(t *testing.T) {
	s := &service{}

	err := s.Init(context.Background())

	assert.NoError(t, err, "Init must not return error")
}

func TestApplyRoutes(t *testing.T) {
	tests := []struct {
		path         string
		mockBehavior func(m *mock_services.MockFbRouter)
		expectedErr  error
	}{
		{
			path: "/fighters",
			mockBehavior: func(m *mock_services.MockFbRouter) {
				m.EXPECT().HandleFunc("/fighters", gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			go func() {
				var panicErr error
				defer func() {
					if r := recover(); r != nil {
						panicErr = fmt.Errorf("Panic: %v", r)
					}
				}()

				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockRouter := mock_services.NewMockFbRouter(ctrl)
				h := &services.ApiHandler{
					ServiceName: "TestHandler",
					Router:      mockRouter,
				}
				service := New(h)

				tc.mockBehavior(mockRouter)

				service.ApplyRoutes()

				assert.Equal(t, tc.expectedErr, panicErr, fmt.Sprintf("ApplyRoutes with %s, must not panic", tc.path))
			}()
		})
	}
}
