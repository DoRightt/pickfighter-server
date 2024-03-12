package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	mock_repo "projects/fb-server/internal/repo/auth/mocks"
	"projects/fb-server/internal/services"
	mock_logger "projects/fb-server/pkg/logger/mocks"
	"projects/fb-server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetCurrentUser(t *testing.T) {
	tests := []struct {
		name           string
		mockBehavior   func(ctx context.Context, mockRepo *mock_repo.MockFbAuthRepo, mockLogger *mock_logger.MockFbLogger)
		req            *http.Request
		expectedStatus int
	}{
		{
			name: "Success",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mlogger *mock_logger.MockFbLogger) {
				userReq := &model.UserRequest{
					UserId: 1,
				}

				user := &model.User{UserId: 1}

				mrepo.EXPECT().FindUser(ctx, userReq).Return(user, nil)
			},
			req: (func() *http.Request {
				req := httptest.NewRequest("Get", "/example", nil)
				ctx := context.WithValue(req.Context(), model.ContextUserId, int32(1))
				req = req.WithContext(ctx)

				return req
			})(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "FindUser error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mlogger *mock_logger.MockFbLogger) {
				userReq := &model.UserRequest{
					UserId: 1,
				}

				mrepo.EXPECT().FindUser(gomock.Any(), userReq).Return(nil, errors.New("error"))
				mlogger.EXPECT().Errorf("Failed to get current user: %s", errors.New("error"))
			},
			req: (func() *http.Request {
				req := httptest.NewRequest("Get", "/example", nil)
				ctx := context.WithValue(req.Context(), model.ContextUserId, int32(1))
				req = req.WithContext(ctx)

				return req
			})(),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repo.NewMockFbAuthRepo(ctrl)
			mockLogger := mock_logger.NewMockFbLogger(ctrl)

			handler := &services.ApiHandler{
				Logger: mockLogger,
			}

			service := &service{
				Repo:       mockRepo,
				ApiHandler: handler,
			}

			w := httptest.NewRecorder()

			tc.mockBehavior(tc.req.Context(), mockRepo, mockLogger)

			service.GetCurrentUser(w, tc.req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
