package repo

import (
	"context"
	"fightbettr.com/fb-server/pkg/model"
	"time"

	"github.com/jackc/pgx/v5"
)

// TxCreateUser creates a new user in the 'fb_users' table.
// If the transaction (tx) is provided, it executes the query within the transaction;
// otherwise, it uses the repository's connection pool to execute the query.
// The user's role (claim) is optional and can be nil if not specified.
// The method returns the newly created user's ID and an error, if any.
func (r *AuthRepo) TxCreateUser(ctx context.Context, tx pgx.Tx, u model.User) (int32, error) {
	query := `INSERT INTO public.fb_users 
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
		if err := tx.QueryRow(ctx, query, args...).Scan(&userId); err != nil {
			return 0, r.DebugLogSqlErr(query, err)
		}
	} else {
		if err := r.GetPool().QueryRow(ctx, query, args...).Scan(&userId); err != nil {
			return 0, r.DebugLogSqlErr(query, err)
		}
	}

	return userId, nil
}
