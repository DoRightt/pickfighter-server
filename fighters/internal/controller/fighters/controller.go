package fighters

import (
	"context"
	"errors"

	"pickfighter.com/fighters/pkg/model"
	logs "pickfighter.com/pkg/logger"
	"pickfighter.com/pkg/pgxs"
	"github.com/jackc/pgx/v5"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type FightersRepository interface {
	pgxs.PickfighterRepo
	SearchFightersCount(ctx context.Context, req *model.FightersRequest) (int32, error)
	SearchFighters(ctx context.Context, req *model.FightersRequest) ([]*model.Fighter, error)
	FindFighter(ctx context.Context, req model.Fighter) (int32, error)
	CreateNewFighter(ctx context.Context, tx pgx.Tx, fighter model.Fighter) (int32, error)
	CreateNewFighterStats(ctx context.Context, tx pgx.Tx, stats model.FighterStats) error
	UpdateFighter(ctx context.Context, tx pgx.Tx, fighter model.Fighter) (int32, error)
	UpdateFighterStats(ctx context.Context, tx pgx.Tx, stats model.FighterStats) error
}

// Controller defines a metadata service controller.
type Controller struct {
	repo FightersRepository
}

// New creates a Fighters service controller.
func New(repo FightersRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}

// SearchFightersCount retrieves the count of fighters based on the provided request.
// It calls the repository's method and returns the count.
// If an error occurs, it logs the error and returns it.
func (c *Controller) SearchFightersCount(ctx context.Context, req *model.FightersRequest) (int32, error) {
	count, err := c.repo.SearchFightersCount(ctx, req)
	if err != nil {
		logs.Errorf("Failed to get fighters count: %s", err)

		// TODO errors package for grpc. Mb it should be handled by a handler on higher level
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.CountFighters, err)

		return 0, err
	}

	return count, nil
}

// SearchFighters retrieves fighters based on the provided request.
// It first retrieves the count of fighters to determine if any exist.
// If no fighters are found, it returns an empty list.
// Otherwise, it calls the repository's method to fetch the fighters and returns them.
// If an error occurs during the process, it logs the error and returns it.
func (c *Controller) SearchFighters(ctx context.Context, req *model.FightersRequest) ([]*model.Fighter, error) {
	count, err := c.repo.SearchFightersCount(ctx, req)
	if err != nil {
		logs.Errorf("Failed to get fighters count: %s", err)

		// TODO errors package for grpc. Mb it should be handled by a handler on higher level
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.CountFighters, err)

		return []*model.Fighter{}, err
	}

	if count == 0 {
		return []*model.Fighter{}, nil
	}

	fighters, err := c.repo.SearchFighters(ctx, req)
	if err != nil {
		logs.Errorf("Failed to find fighters: %s", err)

		// TODO errors package for grpc. Mb it should be handled by a handler on higher level
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.Fighters, err)

		return []*model.Fighter{}, err
	}

	return fighters, nil
}
