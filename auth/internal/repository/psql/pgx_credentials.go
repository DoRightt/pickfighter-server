package psql

import (
	"context"
	"fmt"

	"pickfighter.com/auth/pkg/model"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
)

// TxNewAuthCredentials creates new authentication credentials for a user in the 'fb_user_credentials' table.
// If the transaction (tx) is provided, it executes the query within the transaction;
// otherwise, it uses the repository's connection pool to execute the query.
// The method returns an error if the database operation encounters any issues.
func (r *Repository) TxNewAuthCredentials(ctx context.Context, tx pgx.Tx, uc model.UserCredentials) error {
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
		if _, err := r.GetPool().Exec(ctx, query, args...); err != nil {
			return r.DebugLogSqlErr(query, err)
		}
	}

	return nil
}

// FindUserCredentials retrieves user credentials from the 'fb_user_credentials' table based on the specified request parameters.
// It supports querying by email, token, or user ID. The method constructs a SQL query dynamically based on the provided parameters,
// executes the query using the repository's connection pool, and returns the user credentials if found.
// If the query parameters are invalid or no matching credentials are found, an error is returned.
func (r *Repository) FindUserCredentials(ctx context.Context, req model.UserCredentialsRequest) (model.UserCredentials, error) {
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

	err := r.GetPool().QueryRow(ctx, q).Scan(&c.UserId, &c.Email, &c.Password, &c.Salt, &currentToken, &tokenType, &tokenExpire, &c.Active)
	if err != nil {
		return c, r.DebugLogSqlErr(q, err)
	}

	c.Token = currentToken.String
	c.TokenType = model.TokenType(currentToken.String)
	c.TokenExpire = tokenExpire.Int

	return c, nil
}

// ConfirmCredentialsToken updates the 'fb_user_credentials' table to confirm user credentials based on the provided user ID.
// It sets the 'active' flag to true, updates the token and token type if provided, and sets the token expiration to NULL.
// The method can be used in a transaction (if a non-nil transaction context is provided) or as a standalone operation.
// If the transaction is successful, it confirms the user's credentials; otherwise, it returns an error.
func (r *Repository) ConfirmCredentialsToken(ctx context.Context, tx pgx.Tx, req model.UserCredentialsRequest) error {
	q := `UPDATE public.fb_user_credentials
		SET active = true, token = $2, token_type = $3, token_expire = NULL
		WHERE user_id = $1`

	var token *string
	var tokenType *model.TokenType
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
		if _, err := r.GetPool().Exec(ctx, q, args...); err != nil {
			return r.DebugLogSqlErr(q, err)
		}
	}

	return nil
}

// ResetPassword updates the 'fb_user_credentials' table to reset a user's password based on the provided user credentials.
// It sets the 'active' flag to false, updates the token, token type, and token expiration based on the provided credentials.
// The method is designed to be used when a user requests a password reset.
func (r *Repository) ResetPassword(ctx context.Context, req *model.UserCredentials) error {
	q := `UPDATE public.fb_user_credentials
		SET active = false, token = $2, token_type = $3, token_expire = $4
		WHERE user_id = $1`

	if _, err := r.GetPool().Exec(ctx, q, req.UserId, req.Token, req.TokenType, req.TokenExpire); err != nil {
		return r.DebugLogSqlErr(q, err)
	}

	return nil
}

// UpdatePassword updates the password-related fields of a user in the 'fb_user_credentials' table.
// It is typically used when a user changes their password. The method takes a user's credentials, including the user ID,
// the new hashed password, and the salt used for hashing. The update is performed within a transaction (if provided).
func (r *Repository) UpdatePassword(ctx context.Context, tx pgx.Tx, req model.UserCredentials) error {
	q := `UPDATE public.fb_user_credentials
		SET password_hash = $2, salt = $3
		WHERE user_id = $1`

	if tx != nil {
		if _, err := tx.Exec(ctx, q, req.UserId,
			req.Password, req.Salt); err != nil {
			return r.DebugLogSqlErr(q, err)
		}
	} else {
		if _, err := r.GetPool().Exec(ctx, q, req.UserId,
			req.Password, req.Salt); err != nil {
			return r.DebugLogSqlErr(q, err)
		}
	}

	return nil
}
