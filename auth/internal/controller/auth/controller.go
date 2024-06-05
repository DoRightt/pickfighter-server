package auth

import (
	"context"
	"errors"

	lg "fightbettr.com/auth/pkg/logger"
	"fightbettr.com/auth/pkg/model"
	"fightbettr.com/pkg/pgxs"
	"github.com/jackc/pgx/v5"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type authRepository interface {
	pgxs.FbRepo
	FindUserCredentials(ctx context.Context, req model.UserCredentialsRequest) (model.UserCredentials, error)
	TxNewAuthCredentials(ctx context.Context, tx pgx.Tx, uc model.UserCredentials) error
	ConfirmCredentialsToken(ctx context.Context, tx pgx.Tx, req model.UserCredentialsRequest) error

	ResetPassword(ctx context.Context, req *model.UserCredentials) error
	UpdatePassword(ctx context.Context, tx pgx.Tx, req model.UserCredentials) error
	TxCreateUser(ctx context.Context, tx pgx.Tx, u model.User) (int32, error)

	FindUser(ctx context.Context, req *model.UserRequest) (*model.User, error)
	SearchUsers(ctx context.Context, req *model.UsersRequest) ([]*model.User, error)
	PerformUsersRequestQuery(req *model.UsersRequest) []string
}

// Controller defines a metadata service controller.
type Controller struct {
	repo   authRepository
	Logger lg.FbLogger
}

// New creates a Auth service controller.
func New(repo authRepository) *Controller {
	return &Controller{
		repo:   repo,
		Logger: lg.GetSugared(),
	}
}

func (c *Controller) GracefulShutdown(ctx context.Context, sig string) {
	c.Logger.Warnf("Graceful shutdown. Signal received: %s", sig)
	if c.repo != nil {
		c.repo.GracefulShutdown()
	}
}
