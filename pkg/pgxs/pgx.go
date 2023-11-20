package pgxs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Config struct {
	DataDir  string
	DbUri    string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

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

type Repo struct {
	Logger *zap.SugaredLogger `json:"-" yaml:"-"`
	Pool   *pgxpool.Pool
	Config *Config `json:"-" yaml:"-"`
}

func (db *Repo) GetPoolConfig() (*pgxpool.Config, error) {
	c, err := pgxpool.ParseConfig(db.Config.GetConnString())
	if err != nil {
		return nil, fmt.Errorf("pgxs: unable to parse pgx config: %s", err)
	}

	return c, nil
}

func (db *Repo) GracefulShutdown() {
	if db.Pool != nil {
		db.Pool.Close()
		db.Logger.Infof("Successfully closed postgreSQL connection pool")
	}
}

func (db *Repo) DeleteRecords(ctx context.Context, tableName string) error {
	query := fmt.Sprintf("DELETE FROM %s.%s", "public", tableName)

	_, err := db.Pool.Exec(ctx, query)
	return err
}
