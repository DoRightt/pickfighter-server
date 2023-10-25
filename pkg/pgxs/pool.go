package pgxs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

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

func (db *Repo) ConnectDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	conf, err := db.GetPoolConfig()
	if err != nil {
		return nil, fmt.Errorf("pgxs: Unable to prepare postgres config: %s", err)
	}
	// conf.ConnConfig.TLSConfig = tlsConfig

	return pgxpool.NewWithConfig(ctx, conf)
}
