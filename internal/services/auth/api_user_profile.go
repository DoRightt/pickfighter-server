package auth

import (
	"net/http"
	internalErr "projects/fb-server/pkg/errors"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"
)

func (s *service) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	currentUserId := ctx.Value(model.ContextUserId).(int32)
	user, err := s.Repo.FindUser(ctx, &model.UserRequest{
		UserId: currentUserId,
	})
	if err != nil {
		s.Logger.Errorf("Failed to get current user: %s", err)
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.DBGetUser, err)
		return
	}

	httplib.ResponseJSON(w, model.UserResult{
		User: *user,
	})
}
