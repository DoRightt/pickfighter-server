package common

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	mock_repo "fightbettr.com/fb-server/internal/repo/common/mocks"
	"fightbettr.com/fb-server/internal/services"
	"fightbettr.com/fb-server/pkg/logger"
	"fightbettr.com/fb-server/pkg/model"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetBets(t *testing.T) {
	tests := []struct {
		name           string
		mockBehavior   func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo)
		req            *http.Request
		expectedStatus int
	}{
		{
			name: "No ContextJWTPointer in context",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {

			},
			req:            httptest.NewRequest("GET", "/example", nil),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No UserId in token",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {

			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)
				return createFakeRequestWithToken(token)
			})(),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Success",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchBetsCount(gomock.Any(), int32(1)).Times(1).Return(int32(5), nil)
				mrepo.EXPECT().SearchBets(gomock.Any(), gomock.Any()).Times(1)
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
			name: "SearchBetsCount error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchBetsCount(gomock.Any(), int32(1)).Return(int32(0), errors.New("some error"))
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
			name: "SearchBets error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchBetsCount(gomock.Any(), int32(1)).Return(int32(1), nil)
				mrepo.EXPECT().SearchBets(gomock.Any(), int32(1)).Return(nil, errors.New("some error"))
			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				if err := token.Set(string(model.ContextUserId), float64(1)); err != nil {
					log.Fatalf("Unable to set JWT token userRoles: %s", err)
				}
				return createFakeRequestWithToken(token)
			})(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Bets count is 0",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().SearchBetsCount(gomock.Any(), int32(1)).Return(int32(0), nil)
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

			service.GetBets(w, tc.req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestCreateBet(t *testing.T) {
	tests := []struct {
		name           string
		mockBehavior   func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo)
		req            *http.Request
		expectedStatus int
	}{
		{
			name: "No ContextJWTPointer in context",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {

			},
			req:            httptest.NewRequest("GET", "/example", nil),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No UserId in token",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {

			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				bet := model.Bet{
					BetId:     1,
					FightId:   123,
					UserId:    456,
					FighterId: 789,
				}

				return createFakeRequestWithBody(token, bet)
			})(),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Success",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().CreateBet(gomock.Any(), gomock.Any()).Return(int32(1), nil).Times(1)
			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				if err := token.Set(string(model.ContextUserId), float64(1)); err != nil {
					log.Fatalf("Unable to set JWT token userRoles: %s", err)
				}

				bet := model.Bet{
					BetId:     1,
					FightId:   123,
					UserId:    456,
					FighterId: 789,
				}

				return createFakeRequestWithBody(token, bet)
			})(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "CreateBet Error",
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbCommonRepo) {
				mrepo.EXPECT().CreateBet(gomock.Any(), gomock.Any()).Return(int32(0), errors.New("some error")).Times(1)
			},
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				if err := token.Set(string(model.ContextUserId), float64(1)); err != nil {
					log.Fatalf("Unable to set JWT token userRoles: %s", err)
				}
				return createFakeRequestWithToken(token)
			})(),
			expectedStatus: http.StatusBadRequest,
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

			service.CreateBet(w, tc.req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func createFakeRequestWithToken(token jwt.Token) *http.Request {
	req := httptest.NewRequest("GET", "/example", nil)

	ctx := context.WithValue(req.Context(), model.ContextJWTPointer, token)
	req = req.WithContext(ctx)

	return req
}

func createFakeRequestWithBody(token jwt.Token, body any) *http.Request {
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal("Error while marshal test")
	}

	bodyString := strings.NewReader(string(b))

	req := httptest.NewRequest("POST", "/example", bodyString)

	ctx := context.WithValue(req.Context(), model.ContextJWTPointer, token)
	req = req.WithContext(ctx)

	return req
}

func getFakeToken() (jwt.Token, error) {
	tokenId, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("Unable to generate token id: %s", err)
	}

	token, err := jwt.NewBuilder().
		JwtID(tokenId.String()).
		Issuer("fb-fightbettr").
		Audience([]string{"localhost"}).
		IssuedAt(time.Now()).
		Subject("test").
		Expiration(time.Now().Add(5 * time.Second)).
		Build()

	if err != nil {
		return nil, err
	}

	return token, nil
}
