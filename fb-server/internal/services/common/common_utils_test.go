package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	mock_repo "fightbettr.com/fb-server/internal/repo/common/mocks"
	mock_tx "fightbettr.com/fb-server/internal/repo/mocs"
	"fightbettr.com/fb-server/internal/services"
	internalErr "fightbettr.com/fb-server/pkg/errors"
	"fightbettr.com/fb-server/pkg/httplib"
	"fightbettr.com/fb-server/pkg/logger"
	"fightbettr.com/fb-server/pkg/model"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type TestTx interface {
	pgx.Tx
}

func TestCreateEvent(t *testing.T) {
	tests := []struct {
		name          string
		mockBehavior  func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx)
		eventReq      *model.EventsRequest
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx) {
				mrepo.EXPECT().TxCreateEvent(ctx, mtx, gomock.Any()).Return(int32(1), nil)
			},
			eventReq:      &model.EventsRequest{},
			expectedError: nil,
		},
		{
			name: "error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx) {
				mrepo.EXPECT().TxCreateEvent(ctx, mtx, gomock.Any()).Return(int32(0), errors.New("Error"))
				mtx.EXPECT().Rollback(ctx).Times(1)
			},
			eventReq:      &model.EventsRequest{},
			expectedError: httplib.NewApiErrFromInternalErr(internalErr.New(internalErr.TxUnknown), http.StatusInternalServerError),
		},
		{
			name: "rollback error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx) {
				mrepo.EXPECT().TxCreateEvent(ctx, mtx, gomock.Any()).Return(int32(0), errors.New("Error"))
				mtx.EXPECT().Rollback(gomock.Any()).Return(errors.New("Rollback error")).Times(1)
			},
			eventReq:      &model.EventsRequest{},
			expectedError: httplib.NewApiErrFromInternalErr(internalErr.New(internalErr.TxUnknown), http.StatusInternalServerError),
		},
		{
			name: "success fighter creation",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx) {
				mrepo.EXPECT().TxCreateEvent(ctx, mtx, gomock.Any()).Return(int32(1), nil)
				mrepo.EXPECT().TxCreateEventFight(ctx, mtx, gomock.Any()).Return(nil)
			},
			eventReq: &model.EventsRequest{
				Fights: []model.Fight{
					{
						EventId:       int32(1),
						FighterRedId:  int32(1),
						FighterBlueId: int32(2),
						IsDone:        false,
						IsCanceled:    false,
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "error fighter creation",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx) {
				mrepo.EXPECT().TxCreateEvent(ctx, mtx, gomock.Any()).Return(int32(1), nil)
				mrepo.EXPECT().TxCreateEventFight(ctx, mtx, gomock.Any()).Return(errors.New("Some error"))
				mtx.EXPECT().Rollback(gomock.Any()).Return(errors.New("Rollback error")).Times(1)
			},
			eventReq: &model.EventsRequest{
				Fights: []model.Fight{
					{
						EventId:       int32(1),
						FighterRedId:  int32(1),
						FighterBlueId: int32(2),
						IsDone:        false,
						IsCanceled:    false,
					},
				},
			},
			expectedError: httplib.NewApiErrFromInternalErr(internalErr.New(internalErr.TxUnknown), http.StatusInternalServerError),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockRepo := mock_repo.NewMockFbCommonRepo(ctrl)
			mockTx := mock_tx.NewMockTestTx(ctrl)

			tc.mockBehavior(ctx, mockRepo, mockTx)

			handler := &services.ApiHandler{
				Logger: logger.NewSugared(),
			}
			service := &service{
				Repo:       mockRepo,
				ApiHandler: handler,
			}

			_, err := service.CreateEvent(context.Background(), mockTx, tc.eventReq)

			assert.Equal(t, tc.expectedError, err, fmt.Sprintf("Expected error = %s, but got %s", tc.expectedError, err))
		})
	}
}

func TestCheckEventIsDone(t *testing.T) {
	tests := []struct {
		name          string
		mockBehavior  func(ctx context.Context, m *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx)
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func(ctx context.Context, m *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx) {
				m.EXPECT().GetEventId(ctx, mtx, gomock.Any()).Return(int32(1), nil)
				m.EXPECT().GetUndoneFightsCount(ctx, mtx, int32(1)).Return(0, nil)
				m.EXPECT().SetEventDone(ctx, mtx, int32(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "GetEvent error",
			mockBehavior: func(ctx context.Context, m *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx) {
				m.EXPECT().GetEventId(ctx, mtx, gomock.Any()).Return(int32(0), errors.New("GetEventId error"))
			},
			expectedError: errors.New("GetEventId error"),
		},
		{
			name: "GetUndoneFightsCount error",
			mockBehavior: func(ctx context.Context, m *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx) {
				m.EXPECT().GetEventId(ctx, mtx, gomock.Any()).Return(int32(123), nil)
				m.EXPECT().GetUndoneFightsCount(ctx, mtx, int32(123)).Return(0, errors.New("GetUndoneFightsCount error"))
			},
			expectedError: errors.New("GetUndoneFightsCount error"),
		},
		{
			name: "SetEventDone error",
			mockBehavior: func(ctx context.Context, m *mock_repo.MockFbCommonRepo, mtx *mock_tx.MockTestTx) {
				m.EXPECT().GetEventId(ctx, mtx, gomock.Any()).Return(int32(123), nil)
				m.EXPECT().GetUndoneFightsCount(ctx, mtx, int32(123)).Return(0, nil)
				m.EXPECT().SetEventDone(ctx, mtx, int32(123)).Return(errors.New("SetEventDone error"))
			},
			expectedError: errors.New("SetEventDone error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockRepo := mock_repo.NewMockFbCommonRepo(ctrl)
			mockTx := mock_tx.NewMockTestTx(ctrl)
			service := &service{Repo: mockRepo}

			tc.mockBehavior(ctx, mockRepo, mockTx)

			err := service.CheckEventIsDone(ctx, mockTx, 1)

			assert.Equal(t, tc.expectedError, err, fmt.Sprintf("Expected error = %s, but got %s", tc.expectedError, err))
		})
	}
}

func TestСapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"world", "World"},
		{"", ""},
		{"a", "A"},
		{"Go", "Go"},
		{"123", "123"},
	}

	for _, test := range tests {
		result := capitalize(test.input)

		assert.Equal(t, result, test.expected, "Expected '%s', but got '%s' for input '%s'", test.expected, result, test.input)
	}
}
