package services

import (
	"context"
	"log"
	"path/filepath"
	"projects/fb-server/pkg/logger"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

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
			err := h.Init(tc.context)

			tc.testFunc(h, err)
		})
	}
}

func TestRun(t *testing.T) {
	// TODO
}
