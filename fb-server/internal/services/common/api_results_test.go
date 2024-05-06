package common

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	mock_repo "fightbettr.com/fb-server/internal/repo/common/mocks"
	mock_tx "fightbettr.com/fb-server/internal/repo/mocs"
	"fightbettr.com/fb-server/internal/services"
	mock_logger "fightbettr.com/fb-server/pkg/logger/mocks"
	"fightbettr.com/fb-server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAddResult(t *testing.T) {
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
				mrepo.EXPECT().SetFightResult(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().GetEventId(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().GetUndoneFightsCount(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().SetEventDone(gomock.Any(), gomock.Any(), gomock.Any())

				mtx.EXPECT().Commit(gomock.Any())
			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				bet := model.FightResultRequest{
					FightId:    1,
					WinnerId:   2,
					NotContest: false,
				}

				return createFakeRequestWithBody(token, bet)
			})(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "Empty body, Bad request",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mtx, nil)
				mrepo.EXPECT().SetFightResult(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().GetEventId(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().GetUndoneFightsCount(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().SetEventDone(gomock.Any(), gomock.Any(), gomock.Any())

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
			name: "SetFightResult error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mtx, nil)
				mrepo.EXPECT().SetFightResult(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
				mrepo.EXPECT().GetEventId(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().GetUndoneFightsCount(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().SetEventDone(gomock.Any(), gomock.Any(), gomock.Any())

				mtx.EXPECT().Commit(gomock.Any())
			},
			req:            httptest.NewRequest("POST", "/example", nil),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Tx Commit error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mtx, nil)
				mrepo.EXPECT().SetFightResult(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().GetEventId(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().GetUndoneFightsCount(gomock.Any(), gomock.Any(), gomock.Any())
				mrepo.EXPECT().SetEventDone(gomock.Any(), gomock.Any(), gomock.Any())

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

			service.AddResult(w, tc.req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
