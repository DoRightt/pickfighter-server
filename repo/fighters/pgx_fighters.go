package repo

import (
	"context"
	"projects/fb-server/pkg/model"

	"github.com/jackc/pgx/v5"
)

func (r *FighterRepo) FindFighter(ctx context.Context, req model.Fighter) (int32, error) {
	q := `SELECT fighter_id FROM fb_fighters WHERE name = $1 AND debut_timestamp = $2`
	var fighterId int32

	err := r.Pool.QueryRow(ctx, q, req.Name, req.DebutTimestamp).Scan(&fighterId)
	if err != nil {
		return fighterId, err
	}

	return fighterId, nil
}

func (r *FighterRepo) CreateNewFighter(ctx context.Context, tx pgx.Tx, fighter model.Fighter) (int32, error) {
	qData := `INSERT INTO public.fb_fighters (
		name, nickname, division, status, hometown,
		trains_at, fighting_style, age, height, weight,
		octagon_debut, debut_timestamp, reach, leg_reach, wins,
		loses, draw, fighter_url, image_url
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	RETURNING fighter_id`

	var fighterId int32

	args := []any{
		fighter.Name, fighter.NickName, fighter.Division, fighter.Status, fighter.Hometown,
		fighter.TrainsAt, fighter.FightingStyle, fighter.Age, fighter.Height, fighter.Weight,
		fighter.OctagonDebut, fighter.DebutTimestamp, fighter.Reach, fighter.LegReach, fighter.Wins,
		fighter.Loses, fighter.Draw, fighter.FighterUrl, fighter.ImageUrl,
	}

	if tx != nil {
		if err := tx.QueryRow(ctx, qData, args...).Scan(&fighterId); err != nil {
			return 0, r.DebugLogSqlErr(qData, err)
		}
	} else {
		if err := r.Pool.QueryRow(ctx, qData, args...).Scan(&fighterId); err != nil {
			return 0, r.DebugLogSqlErr(qData, err)
		}
	}

	return fighterId, nil
}

func (r *FighterRepo) CreateNewFighterStats(ctx context.Context, tx pgx.Tx, stats model.FighterStatsReq) error {
	qStats := `INSERT INTO public.fb_fighter_stats (
		fighter_id, total_sig_str_landed, total_sig_str_attempted, str_accuracy, total_tkd_landed, 
		total_tkd_attempted, tkd_accuracy, sig_str_landed, sig_str_absorbed, sig_str_defense,
		takedown_defense, takedown_avg, submission_avg, knockdown_avg, avg_fight_time,
		win_by_ko, win_by_sub, win_by_dec
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`

	args := []any{
		stats.FighterId, stats.TotalSigStrLanded, stats.TotalSigStrAttempted, stats.StrAccuracy, stats.TotalTkdLanded,
		stats.TotalTkdAttempted, stats.TkdAccuracy, stats.SigStrLanded, stats.SigStrAbs, stats.SigStrDefense,
		stats.TakedownDefense, stats.TakedownAvg, stats.SubmissionAvg, stats.KnockdownAvg, stats.AvgFightTime,
		stats.WinByKO, stats.WinBySub, stats.WinByDec,
	}

	if tx != nil {
		_, err := tx.Exec(ctx, qStats, args...)
		if err != nil {
			return err
		}
	} else {
		_, err := r.Pool.Exec(ctx, qStats, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *FighterRepo) UpdateFighter(ctx context.Context, tx pgx.Tx, fighter model.FighterReq) (int32, error) {
	qData := `UPDATE public.fb_fighters SET
		nickname = $2, division = $3, status = $4, hometown = $5, trains_at = $6, 
		fighting_style = $7, age = $8, height = $9, weight = $10, octagon_debut = $11, 
		debut_timestamp = $12, reach = $13, leg_reach = $14, wins = $15, loses = $16,
		draw = $17, fighter_url = $18, image_url = $19
		WHERE fighter_id = $1
		RETURNING fighter_id`

	var fighterId int32

	args := []any{
		fighter.FighterId,
		fighter.NickName, fighter.Division, fighter.Status, fighter.Hometown, fighter.TrainsAt,
		fighter.FightingStyle, fighter.Age, fighter.Height, fighter.Weight, fighter.OctagonDebut,
		fighter.DebutTimestamp, fighter.Reach, fighter.LegReach, fighter.Wins, fighter.Loses,
		fighter.Draw, fighter.FighterUrl, fighter.ImageUrl,
	}

	if tx != nil {
		if err := tx.QueryRow(ctx, qData, args...).Scan(&fighterId); err != nil {
			return 0, r.DebugLogSqlErr(qData, err)
		}
	} else {
		if err := r.Pool.QueryRow(ctx, qData, args...).Scan(&fighterId); err != nil {
			return 0, r.DebugLogSqlErr(qData, err)
		}
	}

	return fighterId, nil
}

func (r *FighterRepo) UpdateFighterStats(ctx context.Context, tx pgx.Tx, stats model.FighterStatsReq) error {
	qStats := `UPDATE public.fb_fighter_stats SET
		total_sig_str_landed = $2, total_sig_str_attempted = $3, str_accuracy = $4, total_tkd_landed = $5, total_tkd_attempted = $6, 
		tkd_accuracy = $7, sig_str_landed = $8, sig_str_absorbed = $9, sig_str_defense = $10, takedown_defense = $11, 
		takedown_avg = $12, submission_avg = $13, knockdown_avg = $14, avg_fight_time = $15, win_by_ko = $16, 
		win_by_sub = $17, win_by_dec = $18
		WHERE fighter_id = $1`

	args := []any{
		stats.FighterId,
		stats.TotalSigStrLanded, stats.TotalSigStrAttempted, stats.StrAccuracy, stats.TotalTkdLanded, stats.TotalTkdAttempted,
		stats.TkdAccuracy, stats.SigStrLanded, stats.SigStrAbs, stats.SigStrDefense, stats.TakedownDefense,
		stats.TakedownAvg, stats.SubmissionAvg, stats.KnockdownAvg, stats.AvgFightTime, stats.WinByKO,
		stats.WinBySub, stats.WinByDec,
	}

	if tx != nil {
		_, err := tx.Exec(ctx, qStats, args...)
		if err != nil {
			return err
		}
	} else {
		_, err := r.Pool.Exec(ctx, qStats, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
