package repo

import (
	"context"
	"fmt"
	"projects/fb-server/pkg/model"
	"strings"
)

func (r *CommonRepo) SearchFightersCount(ctx context.Context, req *model.FightersRequest) (int32, error) {
	q := `SELECT count(*) FROM public.fb_fighters AS f`

	args := r.performFightersQuery(req)
	if len(args) > 0 {
		q += ` WHERE `
		q += strings.Join(args, ` AND `)
	}

	var count int32
	if err := r.Pool.QueryRow(ctx, q).Scan(&count); err != nil {
		return 0, r.DebugLogSqlErr(q, err)
	}

	return count, nil
}

func (r *CommonRepo) SearchFighters(ctx context.Context, req *model.FightersRequest) ([]*model.Fighter, error) {
	q := `SELECT f.fighter_id, f.name, f.nickname, f.division, f.status,
		f.hometown, f.trains_at, f.fighting_style, f.age, f.height,
		f.weight, f.octagon_debut, f.debut_timestamp, f.reach, f.leg_reach,
		f.fighter_url, f.image_url, f.wins, f.loses, f.draw,
		fs.total_sig_str_landed, fs.total_sig_str_attempted, fs.str_accuracy, fs.total_tkd_landed, fs.total_tkd_attempted,
		fs.tkd_accuracy, fs.sig_str_landed, fs.sig_str_absorbed, fs.sig_str_defense, fs.takedown_defense,
		fs.takedown_avg, fs.submission_avg, fs.knockdown_avg, fs.avg_fight_time, fs.win_by_ko,
		fs.win_by_sub, fs.win_by_dec
		FROM public.fb_fighters AS f
		LEFT JOIN public.fb_fighter_stats AS fs ON f.fighter_id = fs.fighter_id`

	args := r.performFightersQuery(req)
	if len(args) > 0 {
		q += ` WHERE `
		q += strings.Join(args, ` AND `)
	}

	rows, err := r.Pool.Query(ctx, q)
	if err != nil {
		return nil, r.DebugLogSqlErr(q, err)
	}
	defer rows.Close()

	var results []*model.Fighter

	for rows.Next() {
		var f model.Fighter
		var fs model.FighterStats

		if err := rows.Scan(
			&f.FighterId, &f.Name, &f.NickName, &f.Division, &f.Status,
			&f.Hometown, &f.TrainsAt, &f.FightingStyle, &f.Age, &f.Height,
			&f.Weight, &f.OctagonDebut, &f.DebutTimestamp, &f.Reach, &f.LegReach,
			&f.FighterUrl, &f.ImageUrl, &f.Wins, &f.Loses, &f.Draw,
			&fs.TotalSigStrLanded, &fs.TotalSigStrAttempted, &fs.StrAccuracy, &fs.TotalTkdLanded, &fs.TotalTkdAttempted,
			&fs.TkdAccuracy, &fs.SigStrLanded, &fs.SigStrAbs, &fs.SigStrDefense, &fs.TakedownDefense,
			&fs.TakedownAvg, &fs.SubmissionAvg, &fs.KnockdownAvg, &fs.AvgFightTime, &fs.WinByKO, &fs.WinBySub, &fs.WinByDec,
		); err != nil {
			return nil, err
		}

		f.Stats = fs

		results = append(results, &f)
	}

	return results, nil
}

func (r *CommonRepo) performFightersQuery(req *model.FightersRequest) []string {
	var args []string
	if req == nil {
		return args
	}

	if req.Status != "" {
		args = append(args, fmt.Sprintf(`f.status = '%s'`, req.Status))
	}

	return args
}
