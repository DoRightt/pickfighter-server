package repo

import (
	"context"
	"projects/fb-server/pkg/model"

	"github.com/jackc/pgx/v5"
)

func (r *CommonRepo) TxCreateEvent(ctx context.Context, tx pgx.Tx, e *model.EventsRequest) (int32, error) {
	q := `INSERT INTO public.fb_events 
	(name)
	VALUES ($1)
	RETURNING event_id`

	args := []any{
		e.Name,
	}

	var eventId int32
	if tx != nil {
		if err := tx.QueryRow(ctx, q, args...).Scan(&eventId); err != nil {
			return 0, r.DebugLogSqlErr(q, err)
		}
	} else {
		if err := r.Pool.QueryRow(ctx, q, args...).Scan(&eventId); err != nil {
			return 0, r.DebugLogSqlErr(q, err)
		}
	}

	return eventId, nil
}
