package repo

import "projects/fb-server/pkg/pgxs"

type FighterRepo struct {
	*pgxs.Repo
}

func New(repo *pgxs.Repo) *FighterRepo {
	return &FighterRepo{
		Repo: repo,
	}
}