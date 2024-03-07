package repo

import (
	"context"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/pgxs"

	"github.com/jackc/pgx/v5"
)

const sep = ` AND `

// FbAuthRepo an interface for handling authentication-related database operations.
type FbAuthRepo interface {
	pgxs.FbRepo
	FindUserCredentials(ctx context.Context, req model.UserCredentialsRequest) (model.UserCredentials, error)
	TxNewAuthCredentials(ctx context.Context, tx pgx.Tx, uc model.UserCredentials) error
	ConfirmCredentialsToken(ctx context.Context, tx pgx.Tx, req model.UserCredentialsRequest) error

	ResetPassword(ctx context.Context, req *model.UserCredentials) error
	UpdatePassword(ctx context.Context, tx pgx.Tx, req model.UserCredentials) error
	TxCreateUser(ctx context.Context, tx pgx.Tx, u model.User) (int32, error)
	
	FindUser(ctx context.Context, req *model.UserRequest) (*model.User, error)
	SearchUsers(ctx context.Context, req *model.UsersRequest) ([]*model.User, error)
	// performUsersRequestQuery(req *model.UsersRequest) []string
}

// AuthRepo represents a repository for handling authentication-related database operations.
// It embeds the *pgxs.Repo, which provides the basic PostgreSQL database operations.
type AuthRepo struct {
	pgxs.FbRepo
}

// New creates and returns a new instance of AuthRepo, initialized with the provided *pgxs.Repo.
func New(r pgxs.FbRepo) FbAuthRepo {
	return &AuthRepo{
		FbRepo: r,
	}
}
