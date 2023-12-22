package repo

import (
	"context"
	"projects/fb-server/pkg/model"
)

func (r *CommonRepo) SearchBetsCount(ctx context.Context, userId int32) (int32, error) {
	q := `SELECT COUNT(*) FROM public.fb_bets WHERE user_id = $1`

	var count int32
	if err := r.Pool.QueryRow(ctx, q, userId).Scan(&count); err != nil {
		return 0, r.DebugLogSqlErr(q, err)
	}

	return count, nil

}
 
func (r *CommonRepo) SearchBets(ctx context.Context, userId int32) ([]*model.Bet, error) {
	q := `SELECT 
	bet_id, user_id, fight_id, bet
	FROM public.fb_bets
	WHERE user_id = $1`

	rows, err := r.Pool.Query(ctx, q, userId)
	if err != nil {
		return nil, r.DebugLogSqlErr(q, err)
	}
	defer rows.Close()

	var bets []*model.Bet
	for rows.Next() {
		var bet model.Bet
		if err := rows.Scan(
			&bet.BetId, &bet.UserId, &bet.FightId, &bet.FighterId,
		); err != nil {
			return nil, r.DebugLogSqlErr(q, err)
		}
		bets = append(bets, &bet)
	}

	return bets, nil

}

func (r *CommonRepo) CreateBet(ctx context.Context, bet *model.Bet) (int32, error) {
	q := `INSERT INTO public.fb_bets 
	(user_id, fight_id, bet)
	VALUES ($1, $2, $3)
	RETURNING bet_id`

	var betId int32
	if err := r.Pool.QueryRow(ctx, q, bet.UserId, bet.FightId, bet.FighterId).Scan(&betId); err != nil {
		return 0, r.DebugLogSqlErr(q, err)
	}

	return betId, nil
}
