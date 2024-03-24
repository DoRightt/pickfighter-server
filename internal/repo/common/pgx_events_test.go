package repo

import (
	"context"
	"log"
	"projects/fb-server/pkg/logger"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/pgxs"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/assert.v1"
)

func TestTxCreateEvent(t *testing.T) {
	updateTables()
	defer viper.Reset()

	ctx := context.Background()
	config := getConfig()

	db, err := pgxs.NewPool(context.Background(), logger.NewSugared(), config)
	require.NoError(t, err)

	commonRepo := New(db)

	req := &model.EventsRequest{
		Name:   "testEvent",
		Fights: []model.Fight{},
	}

	tx, err := commonRepo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	require.NoError(t, err)

	userId, err := commonRepo.TxCreateEvent(ctx, tx, req)
	require.NoError(t, err)

	assert.Equal(t, userId, int32(1))

	userId2, err := commonRepo.TxCreateEvent(ctx, nil, req)
	require.NoError(t, err)

	assert.Equal(t, userId2, int32(2))
}

func TestSearchEventsCount(t *testing.T) {
	// TODO
}

func TestSearchEvents(t *testing.T) {
	// TODO
}

func TestGetEventId(t *testing.T) {
	// TODO
}

func TestGetUndoneFightsCount(t *testing.T) {
	// TODO
}

func TestSetEventDone(t *testing.T) {
	// TODO
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
		_, err := authRepo.GetPool().Exec(context.Background(), "TRUNCATE TABLE public.fb_events RESTART IDENTITY;")
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
