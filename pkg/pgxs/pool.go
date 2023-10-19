package pgxs

import (
	"context"

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

	return s, nil
}
