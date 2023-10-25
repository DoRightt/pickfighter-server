package repo

import (
	"context"
	"fmt"
	"projects/fb-server/pkg/model"
	"time"

	"github.com/jackc/pgx/v5"
)

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
			fmt.Println("ERROR", "1")
			return 0, r.DebugLogSqlErr(query, err)
		}
	} else {
		if err := r.Pool.QueryRow(ctx, query, args...).Scan(&userId); err != nil {
			fmt.Println("ERROR", "2")
			return 0, r.DebugLogSqlErr(query, err)
		}
	}

	return userId, nil
}
