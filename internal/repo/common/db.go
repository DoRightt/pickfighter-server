package repo

import "projects/fb-server/pkg/pgxs"

const sep = ` AND `

// CommonRepo represents a repository for for interacting with fight-related data in the database.
// It embeds the pgxs.Repo, which provides the basic PostgreSQL database operations.
type CommonRepo struct {
	pgxs.FbRepo
}

// New creates and returns a new instance of CommonRepo using the provided pgxs.Repo.
func New(r pgxs.FbRepo) *CommonRepo {
	return &CommonRepo{
		FbRepo: r,
	}
}
