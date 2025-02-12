package psql

import (
	"context"

	eventmodel "github.com/DoRightt/pickfighter-server/events/pkg/model"
	"github.com/jackc/pgx/v5"
)

// TxCreateEventFight creates a new fight in the 'fights' table within a transaction.
// It takes a context, a transaction, and a Fight model.
// It returns an error if the insertion fails.
func (r *Repository) TxCreateEventFight(ctx context.Context, tx pgx.Tx, f eventmodel.Fight) error {
	q := `INSERT INTO
		events.fights(event_id, fighter_red_id, fighter_blue_id)
		VALUES ($1, $2, $3)`

	args := []any{
		f.EventId, f.FighterRedId, f.FighterBlueId,
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

// CreateFightResult creates result row in fight_results table.
// It takes a context, a transaction, and a FightResultRequest.
// It returns an error if the update fails.
func (r *Repository) CreateFightResult(ctx context.Context, tx pgx.Tx, req *eventmodel.FightResultRequest) error {
	q := `INSERT INTO events.fight_results
	(fight_id, winner_id, not_contest, is_draw)
	VALUES ($1, $2, $3, $4);`

	args := []any{
		req.FightId, req.WinnerId, req.NotContest, req.IsDraw,
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

// SetFightResult set is_done field as true of a fight in the 'fights' table.
// It takes a context, a transaction, and a fight id.
// It returns an error if the update fails.
func (r *Repository) SetFightIsDone(ctx context.Context, tx pgx.Tx, fightId int) error {
	q := `UPDATE events.fights
	SET is_done = true
	WHERE fight_id = $1;`

	args := []any{
		fightId,
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
