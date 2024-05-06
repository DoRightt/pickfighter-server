package pgxs

import (
	"context"
	"fightbettr.com/fb-server/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewPool(t *testing.T) {
	lg := logger.NewSugared()

	tests := []struct {
		Name          string
		Context       context.Context
		Logger        *zap.SugaredLogger
		Config        *Config
		ExpectedError error
	}{
		{
			Name:          "Valid configuration",
			Context:       context.Background(),
			Logger:        lg,
			Config:        &Config{},
			ExpectedError: nil,
		},
		{
			Name:          "Empty Config",
			Context:       context.Background(),
			Logger:        lg,
			Config:        nil,
			ExpectedError: ErrEmptyConfig,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			repo, err := NewPool(tc.Context, tc.Logger, tc.Config)

			assert.Equal(t, tc.ExpectedError, err, "Unexpected error")

			if err == nil {
				assert.NotNil(t, repo, "Repo should not be nil")
				assert.NotNil(t, repo.GetPool(), "Pool should not be nil")
			}
		})
	}
}

func TestConntectDBPool(t *testing.T) {
	// TODO
}
