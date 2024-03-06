package common

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	mock_repo "projects/fb-server/internal/repo/common/mocks"
	mock_tx "projects/fb-server/internal/repo/mocs"
	"projects/fb-server/internal/services"
	"projects/fb-server/pkg/logger"
	mock_logger "projects/fb-server/pkg/logger/mocks"
	"projects/fb-server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandleNewEvent(t *testing.T) {
	tests := []struct {
		name           string
		mockBehavior   func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger)
		req            *http.Request
		expectedStatus int
	}{
		{
			name: "Success",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mtx, nil)
				mrepo.EXPECT().TxCreateEvent(gomock.Any(), gomock.Any(), gomock.Any())
				mtx.EXPECT().Commit(gomock.Any())
			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				bet := model.EventsRequest{
					Name:   "Test",
					Fights: []model.Fight{},
				}

				return createFakeRequestWithBody(token, bet)
			})(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "Empty body, Bad request",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mtx, nil)
				mrepo.EXPECT().TxCreateEvent(gomock.Any(), gomock.Any(), gomock.Any())
				mtx.EXPECT().Commit(gomock.Any())
			},
			req:            httptest.NewRequest("POST", "/example", nil),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Begin Tx error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				mlogger.EXPECT().Errorf("Unable to begin transaction: %s", errors.New("error"))
			},
			req:            httptest.NewRequest("POST", "/example", nil),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Tx Commit error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mtx, nil)
				mrepo.EXPECT().TxCreateEvent(gomock.Any(), gomock.Any(), gomock.Any())
				mtx.EXPECT().Commit(gomock.Any()).Return(errors.New("error"))
				mlogger.EXPECT().Errorf("Unable to commit transaction: %s", errors.New("error"))
			},
			req:            httptest.NewRequest("POST", "/example", nil),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockRepo := mock_repo.NewMockFbCommonRepo(ctrl)
			mockLogger := mock_logger.NewMockFbLogger(ctrl)
			mockTx := mock_tx.NewMockTestTx(ctrl)
			handler := &services.ApiHandler{
				Logger: mockLogger,
			}

			service := &service{
				Repo:       mockRepo,
				ApiHandler: handler,
			}

			w := httptest.NewRecorder()

			tc.mockBehavior(ctx, mockRepo, mockTx, mockLogger)

			service.HandleNewEvent(w, tc.req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestGetEvents(t *testing.T) {
	tests := []struct {
		name           string
		mockBehavior   func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo)
		req            *http.Request
		expectedStatus int
	}{
		{
			name: "Success",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchEventsCount(gomock.Any()).Times(1).Return(int32(5), nil)
				mrepo.EXPECT().SearchEvents(gomock.Any()).Times(1)
			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				if err := token.Set(string(model.ContextUserId), float64(1)); err != nil {
					log.Fatalf("Unable to set JWT token userRoles: %s", err)
				}
				return createFakeRequestWithToken(token)
			})(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "SearchEventsCount error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchEventsCount(gomock.Any()).Return(int32(0), errors.New("some error")).AnyTimes()
			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				if err := token.Set(string(model.ContextUserId), float64(1)); err != nil {
					log.Fatalf("Unable to set JWT token userRoles: %s", err)
				}
				return createFakeRequestWithToken(token)
			})(),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Events count is 0",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchEventsCount(gomock.Any()).Return(int32(0), nil)
			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				if err := token.Set(string(model.ContextUserId), float64(1)); err != nil {
					log.Fatalf("Unable to set JWT token userRoles: %s", err)
				}
				return createFakeRequestWithToken(token)
			})(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "SearchEvents error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchEventsCount(gomock.Any()).Return(int32(1), nil)
				mrepo.EXPECT().SearchEvents(gomock.Any()).Return(nil, errors.New("some error"))
			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				if err := token.Set(string(model.ContextUserId), float64(1)); err != nil {
					log.Fatalf("Unable to set JWT token userRoles: %s", err)
				}
				return createFakeRequestWithToken(token)
			})(),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockRepo := mock_repo.NewMockFbCommonRepo(ctrl)
			handler := &services.ApiHandler{
				Logger: logger.NewSugared(),
			}

			service := &service{
				Repo:       mockRepo,
				ApiHandler: handler,
			}

			w := httptest.NewRecorder()

			tc.mockBehavior(ctx, mockRepo)

			service.GetEvents(w, tc.req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
