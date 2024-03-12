package auth

import (
	"context"
	"errors"
	"net/http"
	mock_repo "projects/fb-server/internal/repo/auth/mocks"
	mock_tx "projects/fb-server/internal/repo/mocs"
	"projects/fb-server/internal/services"
	internalErr "projects/fb-server/pkg/errors"
	"projects/fb-server/pkg/httplib"
	mock_logger "projects/fb-server/pkg/logger/mocks"
	"projects/fb-server/pkg/model"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUserCredentials(t *testing.T) {
	tests := []struct {
		name          string
		registerReq   model.RegisterRequest
		mockBehavior  func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger)
		expectedError error
	}{
		{
			name: "No Email",
			registerReq: model.RegisterRequest{
				Name:     "test",
				Password: "12345",
			},
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
			},
			expectedError: httplib.NewApiErrFromInternalErr(internalErr.New(internalErr.AuthFormEmailInvalid)),
		},
		{
			name: "Bad Email",
			registerReq: model.RegisterRequest{
				Name:     "test",
				Password: "12345",
				Email:    "test123",
			},
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
			},
			expectedError: httplib.NewApiErrFromInternalErr(internalErr.New(internalErr.AuthFormEmailInvalid)),
		},
		{
			name: "Bad Password",
			registerReq: model.RegisterRequest{
				Name:     "test",
				Password: "123",
				Email:    "test@mgmail.com",
			},
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
			},
			expectedError: httplib.NewApiErrFromInternalErr(internalErr.New(internalErr.AuthFormPasswordInvalid)),
		},
		{
			name: "TxCreateUser Error, failed to create user",
			registerReq: model.RegisterRequest{
				Name:     "test",
				Password: "Qwerty123456",
				Email:    "test@mgmail.com",
			},
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().TxCreateUser(gomock.Any(), mtx, gomock.Any()).Return(int32(0), errors.New("error"))
				mtx.EXPECT().Rollback(gomock.Any()).Return(errors.New("error"))
				mlogger.EXPECT().Errorf("Unable to rollback transaction: %s", errors.New("error"))
				mlogger.EXPECT().Errorf("Failed to create user during registration transaction: %s", errors.New("error"))
			},
			expectedError: httplib.NewApiErrFromInternalErr(internalErr.New(internalErr.TxUnknown), http.StatusInternalServerError),
		},
		{
			name: "TxCreateUser Error, tx is no unique",
			registerReq: model.RegisterRequest{
				Name:     "test",
				Password: "Qwerty123456",
				Email:    "test@mgmail.com",
			},
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().TxCreateUser(gomock.Any(), mtx, gomock.Any()).Return(int32(0), &pgconn.PgError{Message: "Error", Code: pgerrcode.UniqueViolation})
				mtx.EXPECT().Rollback(gomock.Any()).Return(errors.New("error"))
				mlogger.EXPECT().Errorf("Unable to rollback transaction: %s", errors.New("error"))
			},
			expectedError: httplib.NewApiErrFromInternalErr(internalErr.New(internalErr.TxNotUnique)),
		},
		{
			name: "TxNewAuthCredentials error, failed to create credentials",
			registerReq: model.RegisterRequest{
				Name:     "test",
				Password: "Qwerty123456",
				Email:    "test@mgmail.com",
			},
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().TxCreateUser(gomock.Any(), mtx, gomock.Any()).Return(int32(1), nil)
				mrepo.EXPECT().TxNewAuthCredentials(gomock.Any(), mtx, gomock.Any()).Return(errors.New("error"))
				mtx.EXPECT().Rollback(gomock.Any()).Return(errors.New("error"))
				mlogger.EXPECT().Errorf("Unable to rollback transaction: %s", errors.New("error"))
				mlogger.EXPECT().Errorf("Failed to create user during registration transaction: %s", errors.New("error"))
			},
			expectedError: httplib.NewApiErrFromInternalErr(internalErr.New(internalErr.TxUnknown), http.StatusInternalServerError),
		},
		{
			name: "TxNewAuthCredentials Error, tx is no unique",
			registerReq: model.RegisterRequest{
				Name:     "test",
				Password: "Qwerty123456",
				Email:    "test@mgmail.com",
			},
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mtx *mock_tx.MockTestTx, mlogger *mock_logger.MockFbLogger) {
				mrepo.EXPECT().TxCreateUser(gomock.Any(), mtx, gomock.Any()).Return(int32(1), nil)
				mrepo.EXPECT().TxNewAuthCredentials(gomock.Any(), mtx, gomock.Any()).Return(&pgconn.PgError{Message: "Error", Code: pgerrcode.UniqueViolation})
				mtx.EXPECT().Rollback(gomock.Any()).Return(errors.New("error"))
				mlogger.EXPECT().Errorf("Unable to rollback transaction: %s", errors.New("error"))
			},
			expectedError: httplib.NewApiErrFromInternalErr(internalErr.New(internalErr.TxNotUnique)),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
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

			_, err := service.createUserCredentials(ctx, mockTx, &tc.registerReq)

			assert.Contains(t, err.Error(), tc.expectedError.Error())
		})
	}

	viper.Reset()
}
