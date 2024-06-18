package auth

import (
	"context"

	internalErr "fightbettr.com/auth/pkg/errors"
	"fightbettr.com/auth/pkg/model"
)

// Profile retrieves the user profile based on the provided UserRequest.
// It fetches the user details from the repository and returns them.
// Returns the user profile on success; otherwise returns an error.
func (c *Controller) Profile(ctx context.Context, req *model.UserRequest) (*model.User, error) {
	user, err := c.repo.FindUser(ctx, req)
	if err != nil {
		c.Logger.Errorf("Failed to get current user: %s", err)
		return nil, internalErr.New(internalErr.DBGetUser, err, 801)
	}

	return user, nil
}
