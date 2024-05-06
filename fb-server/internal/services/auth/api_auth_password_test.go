package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	mock_repo "fightbettr.com/fb-server/internal/repo/auth/mocks"
	mock_tx "fightbettr.com/fb-server/internal/repo/mocs"
	"fightbettr.com/fb-server/internal/services"
	mock_logger "fightbettr.com/fb-server/pkg/logger/mocks"
	"fightbettr.com/fb-server/pkg/model"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestResetPassword(t *testing.T) {
	tests := []struct {
		name           string
		req            *http.Request
		mockBehavior   func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger)
		expectedStatus int
	}{
		{
			name: "Success",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.ResetPasswordRequest{
					Email: "test@gmail.com",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				userCredsReq := model.UserCredentialsRequest{
					Email: "test@gmail.com",
				}
				userReq := &model.UserRequest{
					UserId: 1,
				}
				userCreds := model.UserCredentials{UserId: 1}
				mrepo.EXPECT().FindUserCredentials(gomock.Any(), userCredsReq).Return(userCreds, nil)
				mrepo.EXPECT().FindUser(gomock.Any(), userReq).Return(&model.User{UserId: 1, Name: "Test"}, nil)
				mrepo.EXPECT().BeginTx(gomock.Any(), pgx.TxOptions{
					IsoLevel: pgx.Serializable,
				}).Return(mtx, nil)
				mrepo.EXPECT().ResetPassword(gomock.Any(), gomock.Any()).Return(nil)

				mtx.EXPECT().Commit(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Bad request because of empty body",
			req:  httptest.NewRequest("POST", "/example", nil),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {

			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Email is empty",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.ResetPasswordRequest{
					Email: "",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {

			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Email is empty string",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.ResetPasswordRequest{
					Email: " ",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {

			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "FindUserCredentials error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.ResetPasswordRequest{
					Email: "test@gmail.com",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := errors.New("Error")
				userCredsReq := model.UserCredentialsRequest{
					Email: "test@gmail.com",
				}
				mrepo.EXPECT().FindUserCredentials(gomock.Any(), userCredsReq).Return(model.UserCredentials{}, expectedError)
				mlogger.EXPECT().Errorf("Failed to find user credentials: %s", expectedError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "FindUserCredentials NoRows error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.ResetPasswordRequest{
					Email: "test@gmail.com",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := pgx.ErrNoRows
				userCredsReq := model.UserCredentialsRequest{
					Email: "test@gmail.com",
				}
				mrepo.EXPECT().FindUserCredentials(gomock.Any(), userCredsReq).Return(model.UserCredentials{}, expectedError)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "FindUser error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.ResetPasswordRequest{
					Email: "test@gmail.com",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := errors.New("Error")
				userCredsReq := model.UserCredentialsRequest{
					Email: "test@gmail.com",
				}
				userReq := &model.UserRequest{
					UserId: 1,
				}
				mrepo.EXPECT().FindUserCredentials(gomock.Any(), userCredsReq).Return(model.UserCredentials{UserId: 1}, nil)
				mrepo.EXPECT().FindUser(gomock.Any(), userReq).Return(&model.User{}, expectedError)
				mlogger.EXPECT().Errorf("Failed to find user: %s", expectedError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Tx error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.ResetPasswordRequest{
					Email: "test@gmail.com",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := errors.New("Error")
				userCredsReq := model.UserCredentialsRequest{
					Email: "test@gmail.com",
				}
				userReq := &model.UserRequest{
					UserId: 1,
				}
				mrepo.EXPECT().FindUserCredentials(gomock.Any(), userCredsReq).Return(model.UserCredentials{UserId: 1}, nil)
				mrepo.EXPECT().FindUser(gomock.Any(), userReq).Return(&model.User{UserId: 1, Name: "Test"}, nil)
				mrepo.EXPECT().BeginTx(gomock.Any(), pgx.TxOptions{
					IsoLevel: pgx.Serializable,
				}).Return(nil, expectedError)

				mlogger.EXPECT().Errorf("Failed to create registration transaction: %s", expectedError)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "ResetPassword error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.ResetPasswordRequest{
					Email: "test@gmail.com",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := errors.New("Error")
				userCredsReq := model.UserCredentialsRequest{
					Email: "test@gmail.com",
				}
				userReq := &model.UserRequest{
					UserId: 1,
				}
				userCreds := model.UserCredentials{UserId: 1}
				mrepo.EXPECT().FindUserCredentials(gomock.Any(), userCredsReq).Return(userCreds, nil)
				mrepo.EXPECT().FindUser(gomock.Any(), userReq).Return(&model.User{UserId: 1, Name: "Test"}, nil)
				mrepo.EXPECT().BeginTx(gomock.Any(), pgx.TxOptions{
					IsoLevel: pgx.Serializable,
				}).Return(mtx, nil)
				mrepo.EXPECT().ResetPassword(gomock.Any(), gomock.Any()).Return(expectedError)

				mlogger.EXPECT().Errorf("Failed to reset user credentials: %s", expectedError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Tx commit error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.ResetPasswordRequest{
					Email: "test@gmail.com",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := errors.New("Error")
				userCredsReq := model.UserCredentialsRequest{
					Email: "test@gmail.com",
				}
				userReq := &model.UserRequest{
					UserId: 1,
				}
				userCreds := model.UserCredentials{UserId: 1}
				mrepo.EXPECT().FindUserCredentials(gomock.Any(), userCredsReq).Return(userCreds, nil)
				mrepo.EXPECT().FindUser(gomock.Any(), userReq).Return(&model.User{UserId: 1, Name: "Test"}, nil)
				mrepo.EXPECT().BeginTx(gomock.Any(), pgx.TxOptions{
					IsoLevel: pgx.Serializable,
				}).Return(mtx, nil)
				mrepo.EXPECT().ResetPassword(gomock.Any(), gomock.Any()).Return(nil)

				mtx.EXPECT().Commit(gomock.Any()).Return(expectedError)

				mlogger.EXPECT().Errorf("Failed to commit registration transaction: %s", expectedError)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.req.Context()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repo.NewMockFbAuthRepo(ctrl)
			mockLogger := mock_logger.NewMockFbLogger(ctrl)
			mockTx := mock_tx.NewMockTestTx(ctrl)
			handler := &services.ApiHandler{
				Logger: mockLogger,
			}
			service := &service{
				Repo:       mockRepo,
				ApiHandler: handler,
			}

			tc.mockBehavior(ctx, mockRepo, mockTx, mockLogger)

			w := httptest.NewRecorder()

			service.ResetPassword(w, tc.req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestRecoverPassword(t *testing.T) {
	tests := []struct {
		name           string
		req            *http.Request
		mockBehavior   func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger)
		expectedStatus int
	}{
		{
			name: "Success",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token:           "asdasqqwe1243235654ytrewq",
					Password:        "test123",
					ConfirmPassword: "test123",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().FindUserCredentials(gomock.Any(), gomock.Any()).Return(model.UserCredentials{UserId: 1}, nil)
				mrepo.EXPECT().BeginTx(gomock.Any(), pgx.TxOptions{
					IsoLevel: pgx.Serializable,
				}).Return(mtx, nil)
				mrepo.EXPECT().ConfirmCredentialsToken(gomock.Any(), mtx, model.UserCredentialsRequest{
					UserId: 1,
				}).Return(nil)
				mrepo.EXPECT().UpdatePassword(gomock.Any(), mtx, gomock.Any()).Return(nil)

				mtx.EXPECT().Commit(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Bad request because of empty body",
			req:  httptest.NewRequest("POST", "/example", nil),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {

			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Token length is less than 2",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token: "a",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {

			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Token is empty string",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token: " ",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {

			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No password",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token:           "asdasqqwe1243235654ytrewq",
					Password:        "",
					ConfirmPassword: "123456",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {

			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No confirm password",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token:           "asdasqqwe1243235654ytrewq",
					Password:        "123456",
					ConfirmPassword: "",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {

			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No password and no confirm password",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token:           "asdasqqwe1243235654ytrewq",
					Password:        "",
					ConfirmPassword: "",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {

			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "FindUserCredentials error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token:           "asdasqqwe1243235654ytrewq",
					Password:        "test123",
					ConfirmPassword: "test123",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := errors.New("Error")
				mrepo.EXPECT().FindUserCredentials(gomock.Any(), gomock.Any()).Return(model.UserCredentials{}, expectedError)
				mlogger.EXPECT().Errorf("Failed to find user credentials: %s", expectedError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "FindUserCredentials NoRows error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token:           "asdasqqwe1243235654ytrewq",
					Password:        "test123",
					ConfirmPassword: "test123",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := pgx.ErrNoRows

				mrepo.EXPECT().FindUserCredentials(gomock.Any(), gomock.Any()).Return(model.UserCredentials{}, expectedError)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Tx error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token:           "asdasqqwe1243235654ytrewq",
					Password:        "test123",
					ConfirmPassword: "test123",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := errors.New("Error")

				mrepo.EXPECT().FindUserCredentials(gomock.Any(), gomock.Any()).Return(model.UserCredentials{UserId: 1}, nil)
				mrepo.EXPECT().BeginTx(gomock.Any(), pgx.TxOptions{
					IsoLevel: pgx.Serializable,
				}).Return(nil, expectedError)

				mlogger.EXPECT().Errorf("Failed to create registration transaction: %s", expectedError)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "ConfirmCredentialsToken error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token:           "asdasqqwe1243235654ytrewq",
					Password:        "test123",
					ConfirmPassword: "test123",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := errors.New("Error")

				mrepo.EXPECT().FindUserCredentials(gomock.Any(), gomock.Any()).Return(model.UserCredentials{UserId: 1}, nil)
				mrepo.EXPECT().BeginTx(gomock.Any(), pgx.TxOptions{
					IsoLevel: pgx.Serializable,
				}).Return(mtx, nil)
				mrepo.EXPECT().ConfirmCredentialsToken(gomock.Any(), mtx, model.UserCredentialsRequest{
					UserId: 1,
				}).Return(expectedError)

				mlogger.EXPECT().Errorf("Failed to reset user credentials: %s", expectedError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "UpdatePassword error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token:           "asdasqqwe1243235654ytrewq",
					Password:        "test123",
					ConfirmPassword: "test123",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := errors.New("Error")

				mrepo.EXPECT().FindUserCredentials(gomock.Any(), gomock.Any()).Return(model.UserCredentials{UserId: 1}, nil)
				mrepo.EXPECT().BeginTx(gomock.Any(), pgx.TxOptions{
					IsoLevel: pgx.Serializable,
				}).Return(mtx, nil)
				mrepo.EXPECT().ConfirmCredentialsToken(gomock.Any(), mtx, model.UserCredentialsRequest{
					UserId: 1,
				}).Return(nil)
				mrepo.EXPECT().UpdatePassword(gomock.Any(), mtx, gomock.Any()).Return(expectedError)

				mlogger.EXPECT().Errorf("Failed to update user password: %s", expectedError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Tx commit error",
			req: (func() *http.Request {
				token, err := getFakeToken()
				require.NoError(t, err)

				registerReq := model.RecoverPasswordRequest{
					Token:           "asdasqqwe1243235654ytrewq",
					Password:        "test123",
					ConfirmPassword: "test123",
				}

				return createFakeRequestWithBody(token, registerReq)
			})(),
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				expectedError := errors.New("Error")

				mrepo.EXPECT().FindUserCredentials(gomock.Any(), gomock.Any()).Return(model.UserCredentials{UserId: 1}, nil)
				mrepo.EXPECT().BeginTx(gomock.Any(), pgx.TxOptions{
					IsoLevel: pgx.Serializable,
				}).Return(mtx, nil)
				mrepo.EXPECT().ConfirmCredentialsToken(gomock.Any(), mtx, model.UserCredentialsRequest{
					UserId: 1,
				}).Return(nil)
				mrepo.EXPECT().UpdatePassword(gomock.Any(), mtx, gomock.Any()).Return(nil)

				mtx.EXPECT().Commit(gomock.Any()).Return(expectedError)

				mlogger.EXPECT().Errorf("Failed to commit registration transaction: %s", expectedError)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.req.Context()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repo.NewMockFbAuthRepo(ctrl)
			mockLogger := mock_logger.NewMockFbLogger(ctrl)
			mockTx := mock_tx.NewMockTestTx(ctrl)
			handler := &services.ApiHandler{
				Logger: mockLogger,
			}
			service := &service{
				Repo:       mockRepo,
				ApiHandler: handler,
			}

			tc.mockBehavior(ctx, mockRepo, mockTx, mockLogger)

			w := httptest.NewRecorder()

			service.RecoverPassword(w, tc.req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
