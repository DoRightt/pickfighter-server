package repo

import (
	"context"
	"projects/fb-server/pkg/model"

	"github.com/jackc/pgx/v5"
)

func (r *AuthRepo) TxCreateUser(ctx context.Context, tx pgx.Tx, user model.User) (int32, error) {
	return 12, nil
}

func (r *AuthRepo) TxNewAuthCredentials(ctx context.Context, tx pgx.Tx, user model.UserCredentials) error {
	return nil
}
