package pickfighter

import (
	"context"

	authmodel "pickfighter.com/auth/pkg/model"
	eventmodel "pickfighter.com/events/pkg/model"
	fightersmodel "pickfighter.com/fighters/pkg/model"
	"pickfighter.com/pickfighter/pkg/model"
	gatewaymodel "pickfighter.com/pickfighter/pkg/model"
)

type fightersGateway interface {
	SearchFighters(ctx context.Context, req fightersmodel.FightersRequest) ([]*fightersmodel.Fighter, error)
	ServiceHealthCheck() (*model.HealthStatus, error)
}

type authGateway interface {
	Register(ctx context.Context, req *authmodel.RegisterRequest) (*authmodel.UserCredentials, error)
	ConfirmRegistration(ctx context.Context, token string) (bool, error)
	Login(ctx context.Context, req *authmodel.AuthenticateRequest) (*authmodel.AuthenticateResult, error)
	ResetPassword(ctx context.Context, req *authmodel.ResetPasswordRequest) (bool, error)
	PasswordRecover(ctx context.Context, req *authmodel.RecoverPasswordRequest) (bool, error)
	GetCurrentUser(ctx context.Context) (*authmodel.User, error)
	ServiceHealthCheck() (*model.HealthStatus, error)
}

type eventGateway interface {
	CreateEvent(ctx context.Context, req *eventmodel.EventRequest) (*eventmodel.Event, error)
	SearchEvents(ctx context.Context) (*eventmodel.EventsResponse, error)
	CreateBet(ctx context.Context, req *eventmodel.Bet) (*eventmodel.Bet, error)
	SearchBets(ctx context.Context, userId int32) (*eventmodel.BetsResponse, error)
	SetResult(ctx context.Context, req *eventmodel.FightResultRequest) (int32, error)
	ServiceHealthCheck() (*model.HealthStatus, error)
}

// Controller defines a gateway service controller.
type Controller struct {
	authGateway     authGateway
	eventGateway    eventGateway
	fightersGateway fightersGateway
}

// New creates new Controller instance
func New(authGateway authGateway, eventGateway eventGateway, fightersGateway fightersGateway) *Controller {
	return &Controller{
		authGateway,
		eventGateway,
		fightersGateway,
	}
}

// * * * * * Fighters Controller Methods * * * * *

// SearchFighters searches for fighters with the given status using the fightersGateway.
func (c *Controller) SearchFighters(ctx context.Context, req fightersmodel.FightersRequest) ([]*fightersmodel.Fighter, error) {
	fighters, err := c.fightersGateway.SearchFighters(ctx, req)
	if err != nil {
		return nil, err
	}

	return fighters, nil
}

// * * * * * Auth Controller Methods * * * * *

// Register handles the registration of a new user. It takes a context and a
// RegisterRequest, and returns the registered UserCredentials or an error.
func (c *Controller) Register(ctx context.Context, req *authmodel.RegisterRequest) (*authmodel.UserCredentials, error) {
	credentials, err := c.authGateway.Register(ctx, req)
	if err != nil {
		return &authmodel.UserCredentials{}, err
	}

	return credentials, nil
}

// ConfirmRegistration confirms a user's registration with the provided token.
func (c *Controller) ConfirmRegistration(ctx context.Context, token string) (bool, error) {
	ok, err := c.authGateway.ConfirmRegistration(ctx, token)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// // Login authenticates a user with the provided credentials.
func (c *Controller) Login(ctx context.Context, req *authmodel.AuthenticateRequest) (*authmodel.AuthenticateResult, error) {
	token, err := c.authGateway.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// ResetPassword resets a user's password with the provided request details.
func (c *Controller) ResetPassword(ctx context.Context, req *authmodel.ResetPasswordRequest) (bool, error) {
	ok, err := c.authGateway.ResetPassword(ctx, req)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// PasswordRecover initiates the password recovery process for a user with the provided request details.
func (c *Controller) PasswordRecover(ctx context.Context, req *authmodel.RecoverPasswordRequest) (bool, error) {
	ok, err := c.authGateway.PasswordRecover(ctx, req)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// GetCurrentUser retrieves the currently authenticated user.
func (c *Controller) GetCurrentUser(ctx context.Context) (*authmodel.User, error) {
	user, err := c.authGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// * * * * * Events Controller Methods * * * * *

func (c *Controller) CreateEvent(ctx context.Context, req *eventmodel.EventRequest) (*eventmodel.Event, error) {
	event, err := c.eventGateway.CreateEvent(ctx, req)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Controller) SearchEvents(ctx context.Context) (*gatewaymodel.EventsResponse, error) {
	resp, err := c.eventGateway.SearchEvents(ctx)
	if err != nil {
		return nil, err
	}

	fighterIds := c.getFightersIds(resp.Events)

	fReq := fightersmodel.FightersRequest{
		FightersIds: fighterIds,
	}

	fighters, err := c.fightersGateway.SearchFighters(ctx, fReq)
	if err != nil {
		return nil, err
	}

	events := c.eventsPretify(resp.Events, fighters)

	return &gatewaymodel.EventsResponse{Count: resp.Count, Events: events}, nil
}

func (c *Controller) CreateBet(ctx context.Context, req *eventmodel.Bet) (*eventmodel.Bet, error) {
	bet, err := c.eventGateway.CreateBet(ctx, req)
	if err != nil {
		return nil, err
	}

	return bet, nil
}

func (c *Controller) SearchBets(ctx context.Context, userId int32) (*eventmodel.BetsResponse, error) {
	bets, err := c.eventGateway.SearchBets(ctx, userId)
	if err != nil {
		return nil, err
	}

	return bets, nil
}

func (c *Controller) SetResult(ctx context.Context, req *eventmodel.FightResultRequest) (int32, error) {
	id, err := c.eventGateway.SetResult(ctx, req)
	if err != nil {
		return 0, err
	}

	return id, nil
}
