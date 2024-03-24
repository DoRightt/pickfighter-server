package repo

import (
	"context"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/pgxs"

	"github.com/jackc/pgx/v5"
)

type FbFightersRepo interface {
	pgxs.FbRepo
	FindFighter(ctx context.Context, req model.Fighter) (int32, error)
	CreateNewFighter(ctx context.Context, tx pgx.Tx, fighter model.Fighter) (int32, error)
	CreateNewFighterStats(ctx context.Context, tx pgx.Tx, stats model.FighterStats) error
	UpdateFighter(ctx context.Context, tx pgx.Tx, fighter model.Fighter) (int32, error)
	UpdateFighterStats(ctx context.Context, tx pgx.Tx, stats model.FighterStats) error
}

// FighterRepo represents a repository for interacting with fighter-related data in the database.
// It embeds the pgxs.Repo, which provides the basic PostgreSQL database operations.
type FighterRepo struct {
	pgxs.FbRepo
}

// New creates and returns a new instance of FighterRepo using the provided pgxs.Repo.
func New(repo pgxs.FbRepo) *FighterRepo {
	return &FighterRepo{
		FbRepo: repo,
	}
}
