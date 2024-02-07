package repo

import "projects/fb-server/pkg/pgxs"

const sep = ` AND `

type AuthRepo struct {
	*pgxs.Repo
}

func New(r *pgxs.Repo) *AuthRepo {
	return &AuthRepo{
		Repo: r,
	}
}