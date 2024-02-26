package services

import (
	"context"
	"log"
	"path/filepath"
	"projects/fb-server/pkg/cfg"
	"projects/fb-server/pkg/logger"
	"projects/fb-server/pkg/pgxs"
	mock_pgxs "projects/fb-server/pkg/pgxs/mocks"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockService struct {
	*ApiHandler

	Repo *pgxs.Repo `json:"-" yaml:"-"`
}

func (m mockService) Init(ctx context.Context) error {
	return nil
}
func (m mockService) GracefulShutdown()                      {}
func (m mockService) ApplyRoutes()                           {}
func (m mockService) Shutdown(ctx context.Context, s string) {}

func NewMock(h *ApiHandler) ApiService {
	db, err := pgxs.NewPool(context.Background(), logger.NewSugared(), cfg.ViperPostgres())
	if err != nil {
		log.Fatalf("Error connecting to DB: %s\n", err)
	}

	return mockService{
		ApiHandler: h,
		Repo:       db,
	}
}

func initTestConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../../")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
	}
}

func TestNew(t *testing.T) {
	name := "TestName"
	h := New(logger.NewSugared(), name)

	assert.NotNil(t, h, "Handler should not be nil")
	assert.Equal(t, h.ServiceName, name, "Service name should be Test")
}

func TestInit(t *testing.T) {
	initTestConfig()
	certPath := viper.GetString("auth.jwt.cert")
	keyPath := viper.GetString("auth.jwt.key")

	tests := []struct {
		name         string
		context      context.Context
		initSettings func()
		testFunc     func(h *ApiHandler, err error)
	}{
		{
			name:    "OK",
			context: context.Background(),
			initSettings: func() {
				viper.Set("auth.jwt.cert", filepath.Join("..", "..", certPath))
				viper.Set("auth.jwt.key", filepath.Join("..", "..", keyPath))
			},
			testFunc: func(h *ApiHandler, err error) {
				assert.NoError(t, err, "Error should be nil")
				assert.NotNil(t, h.Repo, "Repo should not be nil")
			},
		},
		{
			name:    "Bad certificates path",
			context: context.Background(),
			initSettings: func() {
				viper.Set("auth.jwt.cert", certPath)
				viper.Set("auth.jwt.key", keyPath)
			},
			testFunc: func(h *ApiHandler, err error) {
				assert.NotNil(t, err, "Must be error")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.initSettings()
			h := New(logger.NewSugared(), tc.name)

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_pgxs.NewMockFbRepo(ctl)
			err := h.Init(repo)

			tc.testFunc(h, err)
			viper.Reset()
		})
	}
}

func TestRun(t *testing.T) {
	// TODO
}

func TestAddService(t *testing.T) {
	apiHandler := New(logger.NewSugared(), "TestHandler")
	testService := NewMock(apiHandler)

	apiHandler.AddService("TestService", testService)

	assert.Equal(t, testService, apiHandler.Services["TestService"], "Service should be added correctly")
}

func TestGracefulShutdown(t *testing.T) {
	// TODO
}
