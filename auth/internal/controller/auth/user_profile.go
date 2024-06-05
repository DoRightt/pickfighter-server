package auth

import (
	"context"

	internalErr "fightbettr.com/auth/pkg/errors"
	"fightbettr.com/auth/pkg/model"
)

func (c *Controller) Profile(ctx context.Context, req *model.UserRequest) (*model.User, error) {
	user, err := c.repo.FindUser(ctx, req)
	if err != nil {
		c.Logger.Errorf("Failed to get current user: %s", err)
		return nil, internalErr.New(internalErr.DBGetUser, err, 801)
	}

	return user, nil
}
