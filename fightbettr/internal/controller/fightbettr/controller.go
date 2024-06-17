package fightbettr

import (
	"context"

	authmodel "fightbettr.com/auth/pkg/model"
	fightersmodel "fightbettr.com/fighters/pkg/model"
)

type fightersGateway interface {
	SearchFighters(ctx context.Context, status fightersmodel.FighterStatus) ([]*fightersmodel.Fighter, error)
}

type authGateway interface {
	Register(ctx context.Context, req *authmodel.RegisterRequest) (*authmodel.UserCredentials, error)
	ConfirmRegistration(ctx context.Context, token string) (bool, error)
	Login(ctx context.Context, req *authmodel.AuthenticateRequest) (*authmodel.AuthenticateResult, error)
	ResetPassword(ctx context.Context, req *authmodel.ResetPasswordRequest) (bool, error)
	PasswordRecover(ctx context.Context, req *authmodel.RecoverPasswordRequest) (bool, error)
	GetCurrentUser(ctx context.Context) (*authmodel.User, error)
}

// Controller defines a gateway service controller.
type Controller struct {
	authGateway     authGateway
	fightersGateway fightersGateway
}

func New(authGateway authGateway, fightersGateway fightersGateway) *Controller {
	return &Controller{
		authGateway,
		fightersGateway,
	}
}

// SearchFighters searches for fighters with the given status using the fightersGateway.
func (c *Controller) SearchFighters(ctx context.Context, status string) ([]*fightersmodel.Fighter, error) {
	fighters, err := c.fightersGateway.SearchFighters(ctx, fightersmodel.FighterStatus(status))
	if err != nil {
		return nil, err
	}

	return fighters, nil
}

func (c *Controller) Register(ctx context.Context, req *authmodel.RegisterRequest) (*authmodel.UserCredentials, error) {
	credentials, err := c.authGateway.Register(ctx, req)
	if err != nil {
		return &authmodel.UserCredentials{}, err
	}

	return credentials, nil
}

func (c *Controller) ConfirmRegistration(ctx context.Context, token string) (bool, error) {
	ok, err := c.authGateway.ConfirmRegistration(ctx, token)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (c *Controller) Login(ctx context.Context, req *authmodel.AuthenticateRequest) (*authmodel.AuthenticateResult, error) {
	token, err := c.authGateway.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (c *Controller) ResetPassword(ctx context.Context, req *authmodel.ResetPasswordRequest) (bool, error) {
	ok, err := c.authGateway.ResetPassword(ctx, req)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (c *Controller) PasswordRecover(ctx context.Context, req *authmodel.RecoverPasswordRequest) (bool, error) {
	ok, err := c.authGateway.PasswordRecover(ctx, req)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (c *Controller) GetCurrentUser(ctx context.Context) (*authmodel.User, error) {
	user, err := c.authGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
