package pgxs

import (
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

func (db *Repo) GracefulShutdown() {
	if db.Pool != nil {
		db.Pool.Close()
		db.Logger.Infof("Successfully closed postgreSQL connection pool")
	}
}
