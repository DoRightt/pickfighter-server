package services

import (
	"projects/fb-server/pkg/logger"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadJwtCerts(t *testing.T) {
	tests := []struct {
		name       string
		pathToCert string
		testFunc   func()
	}{
		{
			name: "Correct cert path",
			testFunc: func() {
				h := New(logger.NewSugared(), "TestApp")

				viper.Set("auth.jwt.cert", "../../hack/dev/certs/server-cert.pem")
				viper.Set("auth.jwt.key", "../../hack/dev/certs/server-key.pem")

				err := h.loadJwtCerts()

				assert.NoError(t, err)

				assert.NotNil(t, viper.Get("auth.jwt.signing_key"))
				assert.NotNil(t, viper.Get("auth.jwt.parse_key"))

				viper.Reset()
			},
		},
		{
			name: "Wrong cert path",
			testFunc: func() {
				h := New(logger.NewSugared(), "TestApp")

				viper.Set("auth.jwt.cert", "/wrong/path/to/certs/server-cert.pem")
				viper.Set("auth.jwt.key", "/wrong/path/to/certs/server-key.pem")

				err := h.loadJwtCerts()

				assert.Error(t, err)
				assert.Nil(t, viper.Get("auth.jwt.signing_key"))
				assert.Nil(t, viper.Get("auth.jwt.parse_key"))
			},
		},
		{
			name: "empty cert path",
			testFunc: func() {
				h := New(logger.NewSugared(), "TestApp")

				viper.Set("auth.jwt.cert", "")
				viper.Set("auth.jwt.key", "")

				err := h.loadJwtCerts()

				assert.EqualError(t, err, ErrAuthCertsPathRequired.Error())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc()
		})
	}
}
