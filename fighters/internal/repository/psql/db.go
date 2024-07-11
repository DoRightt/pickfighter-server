package psql

import (
	"context"

	"fightbettr.com/fighters/pkg/cfg"
	"fightbettr.com/pkg/pgxs"
)

// Repository represents a repository for interacting with fighter-related data in the database.
// It embeds the pgxs.Repo, which provides the basic PostgreSQL database operations.
type Repository struct {
	pgxs.FbRepo
}

// New creates and returns a new instance of Fighters Repository using the provided logger
func New(ctx context.Context) (*Repository, error) {
	db, err := pgxs.NewPool(ctx, cfg.ViperPostgres())
	if err != nil {
		return nil, err
	}

	return &Repository{
		FbRepo: db,
	}, nil
}

func (r *Repository) PoolClose() {
	r.GetPool().Close()
}
