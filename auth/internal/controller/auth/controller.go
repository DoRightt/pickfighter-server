package auth

import (
	"context"
	"errors"

	"pickfighter.com/auth/pkg/model"
	"pickfighter.com/pkg/pgxs"
	"github.com/jackc/pgx/v5"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type authRepository interface {
	pgxs.PickfighterRepo
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
	repo authRepository
}

// New creates a Auth service controller.
func New(repo authRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}
