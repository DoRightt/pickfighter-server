package psql

import (
	"context"
	"fmt"
	"strings"

	"pickfighter.com/fighters/pkg/model"
)

// SearchFightersCount retrieves the count of fighters based on the provided FightersRequest.
// It constructs a SQL query to count the number of records in the fb_fighters table, applying
// optional conditions specified in the FightersRequest for filtering. If the request is successful,
// it returns the count of fighters. In case of an error, it returns 0 and the error details.
func (r *Repository) SearchFightersCount(ctx context.Context, req *model.FightersRequest) (int32, error) {
	q := `SELECT count(*) FROM public.fb_fighters AS f`

	args := r.performFightersQuery(req)
	if len(args) > 0 {
		q += ` WHERE `
		q += strings.Join(args, ` AND `)
	}

	var count int32
	if err := r.GetPool().QueryRow(ctx, q).Scan(&count); err != nil {
		return 0, r.DebugLogSqlErr(q, err)
	}

	return count, nil
}

// SearchFighters retrieves a list of fighters based on the provided FightersRequest.
// It constructs a SQL query to join the fb_fighters and fb_fighter_stats tables and applies
// optional conditions specified in the FightersRequest for filtering. The result includes
// information about the fighters and their statistics. If the request is successful, it returns
// a slice of Fighter models. In case of an error, it returns nil and the error details.
func (r *Repository) SearchFighters(ctx context.Context, req *model.FightersRequest) ([]*model.Fighter, error) {
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

	rows, err := r.GetPool().Query(ctx, q)
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

// performFightersQuery constructs the conditions for filtering fighter search based on the provided FightersRequest.
// It returns a slice of string conditions that can be used in the WHERE clause of the SQL query.
// If the provided FightersRequest is nil, an empty slice is returned.
func (r *Repository) performFightersQuery(req *model.FightersRequest) []string {
	var args []string
	if req == nil {
		return args
	}

	if req.Status != "" {
		args = append(args, fmt.Sprintf(`f.status = '%s'`, req.Status))
	}

	if req.FightersIds != nil && len(req.FightersIds) > 0 {
		stringedIds := make([]string, len(req.FightersIds))
		for i, id := range req.FightersIds {
			stringedIds[i] = fmt.Sprintf("%d", id)
		}
		args = append(args, fmt.Sprintf(`f.fighter_id IN (%s)`, strings.Join(stringedIds, ", ")))
	}

	return args
}
