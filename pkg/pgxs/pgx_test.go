package pgxs

import (
	"context"
	"database/sql"
	"log"
	"fightbettr.com/fb-server/pkg/logger"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func initTestConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../../")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
	}
}

func TestGetConnString(t *testing.T) {
	initTestConfig()
	config := &Config{
		Host:     "localhost",
		Port:     "5432",
		Name:     "mydb",
		User:     "user",
		Password: "password",
		DbUri:    "",
	}

	connString := config.GetConnString()

	expectedConnString := "host=localhost port=5432 database=mydb user=user password=password sslmode=disable"
	assert.Equal(t, expectedConnString, connString, "Expected connection string to match")

	config.DbUri = "custom_uri"
	connStringWithUri := config.GetConnString()

	assert.Equal(t, "custom_uri", connStringWithUri, "Expected connection string to match DbUri")
}

func TestGetPoolConfig(t *testing.T) {
	initTestConfig()
	config := &Config{
		Host:     "localhost",
		Port:     "5432",
		Name:     "mydb",
		User:     "user",
		Password: "password",
		DbUri:    "",
	}

	repo := &Repo{
		Config: config,
	}

	poolConfig, err := repo.GetPoolConfig()
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err, "Expected no error from GetPoolConfig")

	expectedPoolConfig := &pgxpool.Config{
		ConnConfig:        poolConfig.ConnConfig,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   30 * time.Minute,
		MaxConns:          12,
		MinConns:          0,
		HealthCheckPeriod: 1 * time.Minute,
	}

	assert.Equal(t, expectedPoolConfig.ConnConfig, poolConfig.ConnConfig, "Expected ConnConfig to match")
	assert.Equal(t, expectedPoolConfig.MaxConnLifetime, poolConfig.MaxConnLifetime, "Expected MaxConnLifetime to match")
	assert.Equal(t, expectedPoolConfig.MaxConnIdleTime, poolConfig.MaxConnIdleTime, "Expected MaxConnIdleTime to match")
	assert.Equal(t, expectedPoolConfig.MaxConns, poolConfig.MaxConns, "Expected MaxConns to match")
	assert.Equal(t, expectedPoolConfig.MinConns, poolConfig.MinConns, "Expected MinConns to match")
	assert.Equal(t, expectedPoolConfig.HealthCheckPeriod, poolConfig.HealthCheckPeriod, "Expected HealthCheckPeriod to match")

	config.DbUri = "custom_uri"
	poolConfigWithDbUri, err := repo.GetPoolConfig()

	assert.Error(t, err, "Expected error")
	assert.Nil(t, poolConfigWithDbUri, "Config with dburi should be nil")

}

func TestGracefulShutdown(t *testing.T) {

	config := &Config{
		DataDir:  viper.GetString("postgres.test.data_dir"),
		DbUri:    viper.GetString("postgres.test.url"),
		Host:     viper.GetString("postgres.test.host"),
		Port:     viper.GetString("postgres.test.port"),
		Name:     viper.GetString("postgres.test.name"),
		User:     viper.GetString("postgres.test.user"),
		Password: viper.GetString("postgres.test.password"),
	}

	db, err := NewPool(context.Background(), logger.NewSugared(), config)
	if err != nil {
		t.Fatalf("Failed to connect to the test database: %v", err)
	}

	db.GracefulShutdown()

	rows, err := db.Pool.Query(context.Background(), "SELECT * FROM table")
	assert.Error(t, err, "Query should return an error because of closed pool")
	rows.Close()
}

func TestDeleteRecords(t *testing.T) {
	config := &Config{
		DataDir:  viper.GetString("postgres.test.data_dir"),
		DbUri:    viper.GetString("postgres.test.url"),
		Host:     viper.GetString("postgres.test.host"),
		Port:     viper.GetString("postgres.test.port"),
		Name:     viper.GetString("postgres.test.name"),
		User:     viper.GetString("postgres.test.user"),
		Password: viper.GetString("postgres.test.password"),
	}

	db, err := NewPool(context.Background(), logger.NewSugared(), config)
	if err != nil {
		t.Fatalf("Failed to connect to the test database: %v", err)
	}
	dbconf, err := db.GetPoolConfig()
	if err != nil {
		t.Fatalf("Failed to get pool config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), dbconf)
	if err != nil {
		t.Fatalf("Failed to connect to the test database: %v", err)
	}
	defer pool.Close()

	_, err = pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS test_table (id serial PRIMARY KEY, name VARCHAR);")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	_, err = pool.Exec(context.Background(), "INSERT INTO test_table (name) VALUES ('test_data');")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	repo := &Repo{
		Pool: pool,
	}

	err = repo.DeleteRecords(context.Background(), "test_table")
	if err != nil {
		t.Fatalf("DeleteRecords failed: %v", err)
	}

	var count int
	err = pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM test_table;").Scan(&count)
	if err != sql.ErrNoRows && count != 0 {
		t.Fatalf("DeleteRecords did not delete records from the table")
	}
}

func TestGetLogger(t *testing.T) {
	config := &Config{
		DataDir:  viper.GetString("postgres.main.data_dir"),
		DbUri:    viper.GetString("postgres.main.url"),
		Host:     viper.GetString("postgres.main.host"),
		Port:     viper.GetString("postgres.main.port"),
		Name:     viper.GetString("postgres.main.name"),
		User:     viper.GetString("postgres.main.user"),
		Password: viper.GetString("postgres.main.password"),
	}

	logger := logger.NewSugared()

	db, err := NewPool(context.Background(), logger, config)
	assert.NoError(t, err, "Expected no error from NewPool")

	assert.NotNil(t, db.GetLogger(), "Expected logger to be not nil")
}
