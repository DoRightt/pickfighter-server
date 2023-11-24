package repo

import "projects/fb-server/pkg/pgxs"

const sep = ` AND `

type CommonRepo struct {
	*pgxs.Repo
}

func New(r *pgxs.Repo) *CommonRepo {
	return &CommonRepo{
		Repo: r,
	}
}
