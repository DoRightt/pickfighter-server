package repo

import "projects/fb-server/pkg/pgxs"

type AuthRepo struct {
	*pgxs.Repo
}

func New(r *pgxs.Repo) *AuthRepo {
	return &AuthRepo{
		Repo: r,
	}
}