package cfg

import (
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/go-playground/assert.v1"
)

func TestViperPostgres(t *testing.T) {
	viper.Set("postgres.main.data_dir", "/test/data_dir")
	viper.Set("postgres.main.url", "postgres://localhost/test")
	viper.Set("postgres.main.host", "localhost")
	viper.Set("postgres.main.port", "5432")
	viper.Set("postgres.main.name", "test")
	viper.Set("postgres.main.user", "test_user")
	viper.Set("postgres.main.password", "test_password")

	cfg := ViperPostgres()

	assert.Equal(t, "/test/data_dir", cfg.DataDir)
	assert.Equal(t, "postgres://localhost/test", cfg.DbUri)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "5432", cfg.Port)
	assert.Equal(t, "test", cfg.Name)
	assert.Equal(t, "test_user", cfg.User)
	assert.Equal(t, "test_password", cfg.Password)
}

func TestViperTestPostgres(t *testing.T) {
	viper.Set("postgres.test.data_dir", "/test/data_dir")
	viper.Set("postgres.test.url", "postgres://localhost/test")
	viper.Set("postgres.test.host", "localhost")
	viper.Set("postgres.test.port", "5432")
	viper.Set("postgres.test.name", "test")
	viper.Set("postgres.test.user", "test_user")
	viper.Set("postgres.test.password", "test_password")

	cfg := ViperTestPostgres()

	assert.Equal(t, "/test/data_dir", cfg.DataDir)
	assert.Equal(t, "postgres://localhost/test", cfg.DbUri)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "5432", cfg.Port)
	assert.Equal(t, "test", cfg.Name)
	assert.Equal(t, "test_user", cfg.User)
	assert.Equal(t, "test_password", cfg.Password)
}
