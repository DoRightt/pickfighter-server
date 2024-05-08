package common

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	mock_repo "fightbettr.com/fb-server/internal/repo/common/mocks"
	"fightbettr.com/fb-server/internal/services"
	"fightbettr.com/fb-server/pkg/logger"
	"fightbettr.com/fb-server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSearchFighters(t *testing.T) {
	tests := []struct {
		name           string
		mockBehavior   func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo)
		req            *http.Request
		expectedStatus int
	}{
		{
			name: "Success",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchFightersCount(gomock.Any(), gomock.Any()).Times(1).Return(int32(5), nil)
				mrepo.EXPECT().SearchFighters(gomock.Any(), gomock.Any()).Times(1)
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
			name: "SearchFightersCount error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchFightersCount(gomock.Any(), gomock.Any()).Return(int32(0), errors.New("some error")).AnyTimes()
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
			name: "Fighters count is 0",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchFightersCount(gomock.Any(), gomock.Any()).Return(int32(0), nil)
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
			name: "SearchFighters error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchFightersCount(gomock.Any(), gomock.Any()).Return(int32(1), nil)
				mrepo.EXPECT().SearchFighters(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
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

			service.SearchFighters(w, tc.req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
