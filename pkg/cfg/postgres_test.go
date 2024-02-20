package cfg

import (
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/go-playground/assert.v1"
)

func TestViperPostgres(t *testing.T) {
	viper.Set("postgres.data_dir", "/test/data_dir")
	viper.Set("postgres.url", "postgres://localhost/test")
	viper.Set("postgres.host", "localhost")
	viper.Set("postgres.port", "5432")
	viper.Set("postgres.name", "test")
	viper.Set("postgres.user", "test_user")
	viper.Set("postgres.password", "test_password")

	cfg := ViperPostgres()

	assert.Equal(t, "/test/data_dir", cfg.DataDir)
	assert.Equal(t, "postgres://localhost/test", cfg.DbUri)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "5432", cfg.Port)
	assert.Equal(t, "test", cfg.Name)
	assert.Equal(t, "test_user", cfg.User)
	assert.Equal(t, "test_password", cfg.Password)
}
