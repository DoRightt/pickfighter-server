package common

import (
	"net/http"
	internalErr "projects/fb-server/pkg/errors"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"
)

func (s *service) SearchFighters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := r.FormValue("status")

	req := &model.FightersRequest{
		Status: capitalize(status),
	}

	count, err := s.Repo.SearchFightersCount(ctx, req)
	if err != nil {
		s.Logger.Errorf("Failed to get fighters count: %s", err)
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.CountFighters, err)
		return
	}
	if count == 0 {
		httplib.ResponseJSON(w, httplib.ListResult{})
		return
	}

	fighters, err := s.Repo.SearchFighters(ctx, req)
	if err != nil {
		s.Logger.Errorf("Failed to find fighters: %s", err)
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.Fighters, err)
		return
	}

	httplib.ResponseJSON(w, httplib.ListResult{
		Results: fighters,
		Count:   count,
	})
}
