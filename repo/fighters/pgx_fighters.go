package repo

import (
	"context"
	"projects/fb-server/pkg/model"

	"github.com/jackc/pgx/v5"
)

func (r *FighterRepo) InsertFightersData(ctx context.Context, tx pgx.Tx, fighters []model.Fighter) error {
	qData := `INSERT INTO fb_fighters (
		name, nickname, division, status, hometown,
		trains_at, fighting_style, age, height, weight,
		octagon_debut, debut_timestamp, reach, leg_reach, wins,
		loses, draw, fighter_url, image_url
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	qStats := `INSERT INTO fb_fighter_stats (
		total_sig_str_landed, total_sig_str_attempted, str_accuracy, total_tkd_landed, total_tkd_attempted,
		tkd_accuracy, sig_str_landed, sig_str_absorbed, sig_str_defense, takedown_defense,
		takedown_avg, submission_avg, knockdown_avg, avg_fight_time, win_by_ko,
		win_by_sub, win_by_dec
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	for _, fighter := range fighters {
		if tx != nil {
			_, err := tx.Exec(ctx,
				qData,
				fighter.Name, fighter.NickName, fighter.Division, fighter.Status, fighter.Hometown,
				fighter.TrainsAt, fighter.FightingStyle, fighter.Age, fighter.Height, fighter.Weight,
				fighter.OctagonDebut, fighter.DebutTimestamp, fighter.Reach, fighter.LegReach, fighter.Wins,
				fighter.Loses, fighter.Draw, fighter.FighterUrl, fighter.ImageUrl,
			)
			if err != nil {
				return err
			}
		} else {
			_, err := r.Pool.Exec(ctx,
				qData,
				fighter.Name, fighter.NickName, fighter.Division, fighter.Status, fighter.Hometown,
				fighter.TrainsAt, fighter.FightingStyle, fighter.Age, fighter.Height, fighter.Weight,
				fighter.OctagonDebut, fighter.DebutTimestamp, fighter.Reach, fighter.LegReach, fighter.Wins,
				fighter.Loses, fighter.Draw, fighter.FighterUrl, fighter.ImageUrl,
			)
			if err != nil {
				return err
			}
		}

		if tx != nil {
			_, err := tx.Exec(ctx,
				qStats,
				fighter.Stats.TotalSigStrLanded, fighter.Stats.TotalSigStrAttempted, fighter.Stats.StrAccuracy, fighter.Stats.TotalTkdLanded, fighter.Stats.TotalTkdAttempted,
				fighter.Stats.TkdAccuracy, fighter.Stats.SigStrLanded, fighter.Stats.SigStrAbs, fighter.Stats.SigStrDefense, fighter.Stats.TakedownDefense,
				fighter.Stats.TakedownAvg, fighter.Stats.SubmissionAvg, fighter.Stats.KnockdownAvg, fighter.Stats.AvgFightTime, fighter.Stats.WinByKO,
				fighter.Stats.WinBySub, fighter.Stats.WinByDec,
			)
			if err != nil {
				return err
			}
		} else {
			_, err := r.Pool.Exec(ctx,
				qStats,
				fighter.Stats.TotalSigStrLanded, fighter.Stats.TotalSigStrAttempted, fighter.Stats.StrAccuracy, fighter.Stats.TotalTkdLanded, fighter.Stats.TotalTkdAttempted,
				fighter.Stats.TkdAccuracy, fighter.Stats.SigStrLanded, fighter.Stats.SigStrAbs, fighter.Stats.SigStrDefense, fighter.Stats.TakedownDefense,
				fighter.Stats.TakedownAvg, fighter.Stats.SubmissionAvg, fighter.Stats.KnockdownAvg, fighter.Stats.AvgFightTime, fighter.Stats.WinByKO,
				fighter.Stats.WinBySub, fighter.Stats.WinByDec,
			)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
