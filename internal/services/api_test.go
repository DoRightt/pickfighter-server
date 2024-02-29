package services_test

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"projects/fb-server/internal/services"
	mock_services "projects/fb-server/internal/services/mocks"
	"projects/fb-server/pkg/cfg"
	"projects/fb-server/pkg/logger"
	mock_logger "projects/fb-server/pkg/logger/mocks"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/pgxs"
	mock_pgxs "projects/fb-server/pkg/pgxs/mocks"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockService struct {
	*services.ApiHandler

	Repo *pgxs.Repo `json:"-" yaml:"-"`
}

func (m mockService) Init(ctx context.Context) error {
	return nil
}
func (m mockService) GracefulShutdown()                      {}
func (m mockService) ApplyRoutes()                           {}
func (m mockService) Shutdown(ctx context.Context, s string) {}

func NewMock(h *services.ApiHandler) services.ApiService {
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
	h := services.New(logger.NewSugared(), name)

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
		testFunc     func(h *services.ApiHandler, err error)
	}{
		{
			name:    "OK",
			context: context.Background(),
			initSettings: func() {
				viper.Set("auth.jwt.cert", filepath.Join("..", "..", certPath))
				viper.Set("auth.jwt.key", filepath.Join("..", "..", keyPath))
			},
			testFunc: func(h *services.ApiHandler, err error) {
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
			testFunc: func(h *services.ApiHandler, err error) {
				assert.NotNil(t, err, "Must be error")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.initSettings()
			h := services.New(logger.NewSugared(), tc.name)

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
	ctrl := gomock.NewController(t)

	service := mock_services.NewMockApiService(ctrl)
	app := &services.ApiHandler{
		Services: map[string]services.ApiService{"mockService": service},
		Logger:   logger.NewSugared(),
	}
	ctx := context.Background()

	app.AddService(model.AuthService, service)

	service.EXPECT().Init(ctx).AnyTimes().Return(nil)

	var panicErr error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = fmt.Errorf("Panic: %v", r)
			}
		}()
		app.Run(ctx)
	}()

	assert.NoError(t, panicErr)
}

func TestAddService(t *testing.T) {
	apiHandler := services.New(logger.NewSugared(), "TestHandler")
	testService := NewMock(apiHandler)

	apiHandler.AddService("TestService", testService)

	assert.Equal(t, testService, apiHandler.Services["TestService"], "Service should be added correctly")
}

func TestGracefulShutdown(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockRepo := mock_pgxs.NewMockFbRepo(ctl)
	mockLogger := mock_logger.NewMockFbLogger(ctl)

	app := &services.ApiHandler{
		Repo:   mockRepo,
		Logger: mockLogger,
	}

	mockRepo.EXPECT().GracefulShutdown()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any())

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var panicErr error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = fmt.Errorf("Panic: %v", r)
			}
		}()
		app.GracefulShutdown(ctx, "SIGTERM")
	}()

	time.Sleep(500 * time.Millisecond)

	assert.NoError(t, ctx.Err())
	assert.Error(t, panicErr)
}
