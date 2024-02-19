package repo

import "projects/fb-server/pkg/pgxs"

// FighterRepo represents a repository for interacting with fighter-related data in the database.
// It embeds the pgxs.Repo, which provides the basic PostgreSQL database operations.
type FighterRepo struct {
	*pgxs.Repo
}

// New creates and returns a new instance of FighterRepo using the provided pgxs.Repo.
func New(repo *pgxs.Repo) *FighterRepo {
	return &FighterRepo{
		Repo: repo,
	}
}