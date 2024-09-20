package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"pickfighter.com/auth/pkg/model"
	"pickfighter.com/auth/pkg/version"
	"pickfighter.com/pkg/pgxs"
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

// HealthCheck returns the current health status of the application.
// It includes information such as the app version, start time, uptime,
// and a message indicating the application's health.
func (c *Controller) HealthCheck() *model.HealthStatus {
	return &model.HealthStatus{
		AppDevVersion: version.DevVersion,
		AppName:       version.Name,
		Timestamp:     time.Now().Format(time.RFC1123),
		AppRunDate:    version.RunDate,
		AppTimeAlive:  time.Now().Unix() - version.RunDate,
		Healthy:       true,
		Message:       fmt.Sprintf("[%s]: I'm fine, thanks!", viper.GetString("app.name")),
	}
}
