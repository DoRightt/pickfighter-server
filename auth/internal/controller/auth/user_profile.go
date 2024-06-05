package auth

import (
	"context"

	"fightbettr.com/auth/pkg/model"
)

func (c *Controller) Profile(ctx context.Context, req *model.UserRequest) (*model.User, error) {
	user, err := c.repo.FindUser(ctx, req)
	if err != nil {
		// TODO handle error
		c.Logger.Errorf("Failed to get current user: %s", err)
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.DBGetUser, err)
		return nil, err
	}

	return user, nil
}
