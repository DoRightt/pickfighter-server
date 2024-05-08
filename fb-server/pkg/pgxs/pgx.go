package pgxs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type FbRepo interface {
	GetPoolConfig() (*pgxpool.Config, error)
	GracefulShutdown()
	DeleteRecords(ctx context.Context, tableName string) error
	ConnectDBPool(ctx context.Context) (*pgxpool.Pool, error)
	DebugLogSqlErr(q string, err error) error
	SanitizeString(s string) string
	GetPool() *pgxpool.Pool
	GetLogger() *zap.SugaredLogger
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// Config is structure that stores data for connecting to a database
type Config struct {
	DataDir  string `json:"data_dir" yaml:"data_dir"`
	DbUri    string `json:"db_uri" yaml:"db_uri"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Name     string `json:"name" yaml:"name"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
}

// GetConnString generates and returns a string based on the data in the config
func (c *Config) GetConnString() string {
	if len(c.DbUri) > 0 {
		return c.DbUri
	}

	connString := fmt.Sprintf("host=%s port=%s database=%s user=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.Name,
		c.User,
		c.Password,
	)

	return connString
}

// Repo is structure for interacting with the database, storing a pool, logger and config
type Repo struct {
	Logger *zap.SugaredLogger `json:"-" yaml:"-"`
	Pool   *pgxpool.Pool      `json:"-" yaml:"-"`
	Config *Config            `json:"-" yaml:"-"`
}

// GetPoolConfig retrieves the pgxpool.Config for the PostgreSQL connection pool.
// It parses the connection string from the Repo's configuration (Config) using pgxpool.ParseConfig.
// Returns the pgxpool.Config and an error if parsing fails.
func (db *Repo) GetPoolConfig() (*pgxpool.Config, error) {
	c, err := pgxpool.ParseConfig(db.Config.GetConnString())
	if err != nil {
		return nil, fmt.Errorf("pgxs: unable to parse pgx config: %s", err)
	}

	return c, nil
}

// BeginTx starts a new transaction with the given options.
func (db *Repo) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return db.Pool.BeginTx(ctx, txOptions)
}

// GetPool returns the pgxpool.Pool stored in the Repo.
func (db *Repo) GetPool() *pgxpool.Pool {
	return db.Pool
}

// GetLogger returns the zap.SugaredLogger stored in the Repo.
func (db *Repo) GetLogger() *zap.SugaredLogger {
	return db.Logger
}

// GracefulShutdown checks whether the value is in the pool and if so, closes it and logs the message
func (db *Repo) GracefulShutdown() {
	if db.Pool != nil {
		db.Pool.Close()
		db.Logger.Infof("Successfully closed postgreSQL connection pool")
	}
}

// DeleteRecords deletes records from the table whose name is passed as an argument
func (db *Repo) DeleteRecords(ctx context.Context, tableName string) error {
	if len(tableName) == 0 {
		return fmt.Errorf("pgxs: table name is empty")
	}
	query := fmt.Sprintf("DELETE FROM %s.%s", "public", tableName)

	_, err := db.Pool.Exec(ctx, query)
	return err
}
