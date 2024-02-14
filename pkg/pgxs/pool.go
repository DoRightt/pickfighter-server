package pgxs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// NewPool creates a new Repo with a configured PostgreSQL connection pool.
// It requires a context, a SugaredLogger, and a database configuration (Config).
// If the configuration is nil, it returns an error.
func NewPool(ctx context.Context, lg *zap.SugaredLogger, conf *Config) (*Repo, error) {
	if conf == nil {
		return nil, ErrEmptyConfig
	}

	s := &Repo{
		Logger: lg.Named("pgx_pool"),
		Config: conf,
	}

	pool, err := s.ConnectDBPool(ctx)
	if err != nil {
		return nil, err
	}

	s.Pool = pool

	return s, nil
}

// ConnectDBPool creates and returns a new pgxpool.Pool using the configured context and pgxpool.Config.
// It uses the GetPoolConfig method to obtain the configuration.
// Returns the created pgxpool.Pool and an error if there is any issue.
func (db *Repo) ConnectDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	conf, err := db.GetPoolConfig()
	if err != nil {
		return nil, fmt.Errorf("pgxs: Unable to prepare postgres config: %s", err)
	}
	// conf.ConnConfig.TLSConfig = tlsConfig

	return pgxpool.NewWithConfig(ctx, conf)
}
