package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
	eventmodel "pickfighter.com/events/pkg/model"
)

// SearchBetsCount retrieves the count of bets for a given user ID from the 'pf_bets' table.
// It takes a context and a user ID, and returns the count of bets or an error if the query fails.
func (r *Repository) SearchBetsCount(ctx context.Context, userId int32) (int32, error) {
	q := `SELECT COUNT(*) FROM public.pf_bets WHERE user_id = $1`

	var count int32
	if err := r.GetPool().QueryRow(ctx, q, userId).Scan(&count); err != nil {
		return 0, r.DebugLogSqlErr(q, err)
	}

	return count, nil
}

// SearchBets retrieves a list of bets for a given user ID from the 'pf_bets' table.
// It takes a context and a user ID, and returns a slice of Bet models or an error if the query fails.
func (r *Repository) SearchBets(ctx context.Context, userId int32) ([]*eventmodel.Bet, error) {
	q := `SELECT 
	bet_id, user_id, fight_id, bet
	FROM public.pf_bets
	WHERE user_id = $1`

	rows, err := r.GetPool().Query(ctx, q, userId)
	if err != nil {
		return nil, r.DebugLogSqlErr(q, err)
	}
	defer rows.Close()

	var bets []*eventmodel.Bet
	for rows.Next() {
		var bet eventmodel.Bet
		if err := rows.Scan(
			&bet.BetId, &bet.UserId, &bet.FightId, &bet.FighterId,
		); err != nil {
			return nil, r.DebugLogSqlErr(q, err)
		}
		bets = append(bets, &bet)
	}

	return bets, nil
}

// CreateBet inserts a new bet into the 'pf_bets' table.
// It takes a context, a Bet model and returns the newly created bet's ID
// or an error if the insertion fails.
func (r *Repository) TxCreateBet(ctx context.Context, tx pgx.Tx, bet *eventmodel.Bet) (int32, error) {
	q := `INSERT INTO public.pf_bets 
	(user_id, fight_id, bet)
	VALUES ($1, $2, $3)
	RETURNING bet_id`

	var betId int32
	if tx != nil {
		if err := tx.QueryRow(ctx, q, bet.UserId, bet.FightId, bet.FighterId).Scan(&betId); err != nil {
			return 0, r.DebugLogSqlErr(q, err)
		}
	} else {
		if err := r.GetPool().QueryRow(ctx, q, bet.UserId, bet.FightId, bet.FighterId).Scan(&betId); err != nil {
			return 0, r.DebugLogSqlErr(q, err)
		}
	}

	return betId, nil
}
