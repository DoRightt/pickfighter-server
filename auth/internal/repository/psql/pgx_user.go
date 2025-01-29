package psql

import (
	"context"
	"fmt"
	"time"

	"pickfighter.com/auth/pkg/model"

	"github.com/jackc/pgx/v5"
)

// TxCreateUser creates a new user in the 'users' table.
// If the transaction (tx) is provided, it executes the query within the transaction;
// otherwise, it uses the repository's connection pool to execute the query.
// The user's role (claim) is optional and can be nil if not specified.
// The method returns the newly created user's ID and an error, if any.
func (r *Repository) TxCreateUser(ctx context.Context, tx pgx.Tx, u model.User) (int32, error) {
	q := `INSERT INTO auth.users 
		(name, claim, created_at)
		VALUES ($1, $2, $3)
		RETURNING user_id`

	var role *string

	if len(u.Claim) > 0 {
		role = &u.Claim
	}

	args := []any{
		u.Name, role, time.Now().Unix(),
	}

	var userId int32
	if tx != nil {
		if err := tx.QueryRow(ctx, q, args...).Scan(&userId); err != nil {
			return 0, r.DebugLogSqlErr(q, err)
		}
	} else {
		if err := r.GetPool().QueryRow(ctx, q, args...).Scan(&userId); err != nil {
			return 0, r.DebugLogSqlErr(q, err)
		}
	}

	return userId, nil
}

func (r *Repository) PatchUser(ctx context.Context, tx pgx.Tx, userId int32, param string, val any) error {
	q := fmt.Sprintf(`UPDATE auth.users SET
		updated_at = $2,
		%s = $3
	WHERE user_id = $1`, param)

	args := []any{userId, time.Now().Unix(), val}
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
