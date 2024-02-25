package services

import (
	"projects/fb-server/pkg/logger"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadJwtCerts(t *testing.T) {
	h := New(logger.NewSugared(), "TestApp")

	viper.Set("auth.jwt.cert", "../../hack/dev/certs/server-cert.pem")
	viper.Set("auth.jwt.key", "../../hack/dev/certs/server-key.pem")

	err := h.loadJwtCerts()

	assert.NoError(t, err)

	assert.NotNil(t, viper.Get("auth.jwt.signing_key"))
	assert.NotNil(t, viper.Get("auth.jwt.parse_key"))

	viper.Reset()
}
