package repo

import "projects/fb-server/pkg/pgxs"

const sep = ` AND `

// AuthRepo represents a repository for handling authentication-related database operations.
// It embeds the *pgxs.Repo, which provides the basic PostgreSQL database operations.
type AuthRepo struct {
	*pgxs.Repo
}

// New creates and returns a new instance of AuthRepo, initialized with the provided *pgxs.Repo.
func New(r *pgxs.Repo) *AuthRepo {
	return &AuthRepo{
		Repo: r,
	}
}