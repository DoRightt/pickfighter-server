package repo

import (
	"context"
	"fmt"
	"projects/fb-server/pkg/model"

	"github.com/jackc/pgtype"
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

func (r *AuthRepo) FindUserCredentials(ctx context.Context, req model.UserCredentialsRequest) (model.UserCredentials, error) {
	q := `SELECT user_id, email, password_hash, salt, token, token_type, token_expire, active
		FROM public.fb_user_credentials`

	if req.Email != "" {
		q += ` WHERE email = '` + r.SanitizeString(req.Email) + `'`
	} else if req.Token != "" {
		q += ` WHERE token = '` + r.SanitizeString(req.Token) + `'`
	} else if req.UserId > 0 {
		q += fmt.Sprintf(` WHERE user_id = %d`, req.UserId)
	} else {
		return model.UserCredentials{}, fmt.Errorf("no query parameters")
	}

	var c model.UserCredentials
	var currentToken, tokenType pgtype.Varchar
	var tokenExpire pgtype.Int8

	err := r.Pool.QueryRow(ctx, q).Scan(&c.UserId, &c.Email, &c.Password, &c.Salt, &currentToken, &tokenType, &tokenExpire, &c.Active)
	if err != nil {
		return c, r.DebugLogSqlErr(q, err)
	}

	c.Token = currentToken.String
	c.TokenType = currentToken.String
	c.TokenExpire = tokenExpire.Int

	return c, nil
}

func (r *AuthRepo) ConfirmCredentialsToken(ctx context.Context, tx pgx.Tx, req model.UserCredentialsRequest) error {
	q := `UPDATE public.fb_user_credentials
		SET active = true, token = $2, token_type = $3, token_expire = NULL
		WHERE user_id = $1`

	var token, tokenType *string
	if len(req.Token) > 0 {
		token = &req.Token
	}
	if len(req.TokenType) > 0 {
		tokenType = &req.TokenType
	}

	args := []any{req.UserId, token, tokenType}
	if tx != nil {
		if _, err := tx.Exec(ctx, q, args...); err != nil {
			return r.DebugLogSqlErr(q, err)
		}
	} else {
		if _, err := r.Pool.Exec(ctx, q, args...); err != nil {
			return r.DebugLogSqlErr(q, err)
		}
	}

	return nil
}
