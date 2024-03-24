package repo

import (
	"context"
	"log"
	mock_repo "projects/fb-server/internal/repo/auth/mocks"
	"projects/fb-server/pkg/logger"
	"projects/fb-server/pkg/pgxs"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockFbAuthRepo(ctrl)

	authRepo := New(mockRepo)

	assert.NotNil(t, authRepo, "authRepo should not be nil")

	switch v := authRepo.(type) {
	case FbAuthRepo:
		return
	default:
		t.Errorf("Expected Repo interface is FbAuthRepo, but got %v", v)
	}
}

func updateTables() {
	initTestConfig()

	config := &pgxs.Config{
		DataDir:  viper.GetString("postgres.test.data_dir"),
		DbUri:    viper.GetString("postgres.test.url"),
		Host:     viper.GetString("postgres.test.host"),
		Port:     viper.GetString("postgres.test.port"),
		Name:     viper.GetString("postgres.test.name"),
		User:     viper.GetString("postgres.test.user"),
		Password: viper.GetString("postgres.test.password"),
	}

	db, err := pgxs.NewPool(context.Background(), logger.NewSugared(), config)
	if err != nil {
		log.Fatal(err)
	}

	authRepo := New(db)

	url := viper.GetString("postgres.test.url")
	parts := strings.Split(url, "/")
	dbName := parts[len(parts)-1]
	testName := "test_table"

	if dbName == testName {
		_, err := authRepo.GetPool().Exec(context.Background(), "TRUNCATE TABLE public.fb_users, public.fb_user_credentials RESTART IDENTITY;")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func initTestConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../../../")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
	}
}

func getConfig() *pgxs.Config {
	config := &pgxs.Config{
		DataDir:  viper.GetString("postgres.test.data_dir"),
		DbUri:    viper.GetString("postgres.test.url"),
		Host:     viper.GetString("postgres.test.host"),
		Port:     viper.GetString("postgres.test.port"),
		Name:     viper.GetString("postgres.test.name"),
		User:     viper.GetString("postgres.test.user"),
		Password: viper.GetString("postgres.test.password"),
	}

	return config
}
