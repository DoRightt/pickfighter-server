package auth

import (
	"net/http"
	internalErr "fightbettr.com/fb-server/pkg/errors"
	"fightbettr.com/fb-server/pkg/httplib"
	"fightbettr.com/fb-server/pkg/model"
)

// GetCurrentUser retrieves information about the currently authenticated user.
// It extracts the user ID from the request context, queries the database for the user's details,
// and responds with a JSON representation of the user.
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
