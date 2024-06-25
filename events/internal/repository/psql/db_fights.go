package psql

import (
	"context"

	eventmodel "fightbettr.com/events/pkg/model"
	"github.com/jackc/pgx/v5"
)

// TxCreateEventFight creates a new fight in the 'fb_fights' table within a transaction.
// It takes a context, a transaction, and a Fight model.
// It returns an error if the insertion fails.
func (r *Repository) TxCreateEventFight(ctx context.Context, tx pgx.Tx, f eventmodel.Fight) error {
	q := `INSERT INTO
		public.fb_fights(event_id, fighter_red_id, fighter_blue_id, is_done, not_contest)
		VALUES ($1, $2, $3, $4, $5)`

	args := []any{
		f.EventId, f.FighterRedId, f.FighterBlueId, f.IsDone, f.NotContest,
	}

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

// SetFightResult updates the result of a fight in the 'fb_fights' table.
// It takes a context, a transaction, and a FightResultRequest.
// It returns an error if the update fails.
func (r *Repository) SetFightResult(ctx context.Context, tx pgx.Tx, req *eventmodel.FightResultRequest) error {
	q := `UPDATE fb_fights
	SET result = $1, not_contest = $2, is_done = true
	WHERE fight_id = $3;`

	args := []any{
		req.WinnerId, req.NotContest, req.FightId,
	}

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
