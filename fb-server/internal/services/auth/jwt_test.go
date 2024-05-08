package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	mock_repo "fightbettr.com/fb-server/internal/repo/auth/mocks"
	"fightbettr.com/fb-server/internal/services"
	mock_logger "fightbettr.com/fb-server/pkg/logger/mocks"
	"fightbettr.com/fb-server/pkg/model"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateJWTToken(t *testing.T) {
	tests := []struct {
		name         string
		creds        model.UserCredentials
		authReq      model.AuthenticateRequest
		mockBehavior func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mlogger *mock_logger.MockFbLogger)
		err          bool
	}{
		{
			name: "Success",
			creds: model.UserCredentials{
				UserId: 1,
				Email:  "test@gmail.com",
			},
			authReq: model.AuthenticateRequest{
				Email:   "test@gmail.com",
				Subject: "test",
			},
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mlogger *mock_logger.MockFbLogger) {
				creds := ctx.Value("creds").(*model.UserCredentials)
				req := ctx.Value("authReq").(model.AuthenticateRequest)
				testUser := model.User{
					UserId: creds.UserId,
					Email:  creds.Email,
					Name:   "test",
				}

				mrepo.EXPECT().FindUser(gomock.Any(), &model.UserRequest{
					UserId: creds.UserId,
				}).Return(&testUser, nil)

				mlogger.EXPECT().Debugf("Issuing JWT token for User [%d:%s:%s]", creds.UserId, creds.Email, req.Subject)
			},
			err: false,
		},
		{
			name: "User is not found",
			creds: model.UserCredentials{
				UserId: 1,
				Email:  "test@gmail.com",
			},
			authReq: model.AuthenticateRequest{
				Email:   "test@gmail.com",
				Subject: "test",
			},
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mlogger *mock_logger.MockFbLogger) {
				creds := ctx.Value("creds").(*model.UserCredentials)

				mrepo.EXPECT().FindUser(gomock.Any(), &model.UserRequest{
					UserId: creds.UserId,
				}).Return(nil, pgx.ErrNoRows)

				mlogger.EXPECT().Errorf("Failed to get user: %s", pgx.ErrNoRows)
			},
			err: true,
		},
		{
			name: "JWT signing error",
			creds: model.UserCredentials{
				UserId: 1,
				Email:  "test@gmail.com",
			},
			authReq: model.AuthenticateRequest{
				Email:   "test@gmail.com",
				Subject: "test",
			},
			mockBehavior: func(ctx context.Context, mrepo *mock_repo.MockFbAuthRepo, mlogger *mock_logger.MockFbLogger) {
				viper.Reset()

				creds := ctx.Value("creds").(*model.UserCredentials)
				req := ctx.Value("authReq").(model.AuthenticateRequest)
				testUser := model.User{
					UserId: creds.UserId,
					Email:  creds.Email,
					Name:   "test",
				}

				mrepo.EXPECT().FindUser(gomock.Any(), &model.UserRequest{
					UserId: creds.UserId,
				}).Return(&testUser, nil)

				mlogger.EXPECT().Debugf("Issuing JWT token for User [%d:%s:%s]", creds.UserId, creds.Email, req.Subject)

				mlogger.EXPECT().Errorf("failed to generate signed payload: %s\n", gomock.Any())
			},
			err: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "creds", &tc.creds)
			ctx = context.WithValue(ctx, "authReq", tc.authReq)

			ctrl := gomock.NewController(t)
			mockRepo := mock_repo.NewMockFbAuthRepo(ctrl)
			mockLogger := mock_logger.NewMockFbLogger(ctrl)
			handler := &services.ApiHandler{
				Logger: mockLogger,
			}
			service := &service{
				Repo:       mockRepo,
				ApiHandler: handler,
			}

			viper.Set("auth.jwt.cert", "../../../hack/dev/certs/server-cert.pem")
			viper.Set("auth.jwt.key", "../../../hack/dev/certs/server-key.pem")
			defer viper.Reset()

			loadJwtCerts()

			tc.mockBehavior(ctx, mockRepo, mockLogger)

			_, err := service.createJWTToken(ctx, &tc.creds, tc.authReq)

			if tc.err {
				assert.Error(t, err, "Should be error")
			} else {
				assert.NoError(t, err, "Unexpected error")
			}
		})
	}
}

func loadJwtCerts() error {
	certPath := viper.GetString("auth.jwt.cert")
	keyPath := viper.GetString("auth.jwt.key")

	hasRsaKeys := len(certPath) > 0 && len(keyPath) > 0

	if !hasRsaKeys {
		return errors.New("error")
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		fmt.Println("ERR", err)
		return err
	}

	viper.Set("auth.jwt.signing_key", cert.PrivateKey)

	clientCert, err := os.ReadFile(certPath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(clientCert)
	readCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	viper.Set("auth.jwt.parse_key", readCert.PublicKey)

	return nil
}
