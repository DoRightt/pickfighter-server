package psql

import (
	"context"

	"pickfighter.com/pkg/pgxs"
)

// Repository represents a repository for interacting with fighter-related data in the database.
// It embeds the pgxs.Repo, which provides the basic PostgreSQL database operations.
type Repository struct {
	pgxs.PickfighterRepo
}

// New creates and returns a new instance of Fighters Repository
func New(ctx context.Context, cfg *pgxs.Config) (*Repository, error) {
	db, err := pgxs.NewPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		PickfighterRepo: db,
	}, nil
}

func (r *Repository) PoolClose() {
	r.GetPool().Close()
}
