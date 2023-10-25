package repo

import (
	"context"
	"projects/fb-server/pkg/model"

	"github.com/jackc/pgx/v5"
)

func (r *AuthRepo) TxNewAuthCredentials(ctx context.Context, tx pgx.Tx, uc model.UserCredentials) error {
	query := `INSERT INTO
		public.fb_user_credentials(user_id, email, password_hash, salt, token, token_type, token_expire, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	args := []any{
		uc.UserId, uc.Email, uc.Password,
		uc.Salt, uc.Token, uc.TokenType, uc.TokenExpire,
		uc.Active,
	}

	if tx != nil {
		if _, err := tx.Exec(ctx, query, args...); err != nil {
			return r.DebugLogSqlErr(query, err)
		}
	} else {
		if _, err := r.Pool.Exec(ctx, query, args...); err != nil {
			return r.DebugLogSqlErr(query, err)
		}
	}

	return nil

}
