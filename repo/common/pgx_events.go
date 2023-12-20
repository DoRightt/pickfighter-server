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

func (r *CommonRepo) SearchEventsCount(ctx context.Context) (int32, error) {
	limit := 5

	q := ` WITH ranked_events AS (
		SELECT
			*,
			ROW_NUMBER() OVER (PARTITION BY is_done ORDER BY event_id) AS rn
		FROM public.fb_events
		WHERE is_done IN (false, true)
	),
	filtered_events AS (
		SELECT
			event_id,
			is_done
		FROM ranked_events
		WHERE (is_done = false AND rn <= $1) OR (is_done = true AND rn > (SELECT COUNT(*) FROM public.fb_events WHERE is_done = true) - $1)
	)
	SELECT COUNT(*) FROM filtered_events`

	var count int32
	if err := r.Pool.QueryRow(ctx, q, limit).Scan(&count); err != nil {
		return 0, r.DebugLogSqlErr(q, err)
	}

	return count, nil
}

func (r *CommonRepo) SearchEvents(ctx context.Context) ([]*model.FullEventResponse, error) {
	limit := 5

	q := `WITH ranked_events AS (
		SELECT
			*,
			ROW_NUMBER() OVER (PARTITION BY is_done ORDER BY event_id) AS rn
		FROM public.fb_events
		WHERE is_done IN (false, true)
	),
	filtered_events AS (
		SELECT
			event_id, name, is_done
		FROM ranked_events
		WHERE (is_done = false AND rn <= $1) OR (is_done = true AND rn > (SELECT COUNT(*) FROM public.fb_events WHERE is_done = true) - $1)
	)
	SELECT
		e.event_id, e.name, e.is_done AS is_event_done, f.fight_id, f.is_done AS is_fight_done,
		f.not_contest, f.created_at, f.fight_date, f.result,
		f.fighter_red_id, red.name AS red_fighter_name, red.nickname AS red_fighter_nickname, red.division AS red_fighter_division, red.fighter_url AS red_fighter_url,
		red.image_url AS red_fighter_image, red.wins AS red_fighter_wins, red.loses AS red_fighter_loses, red.draw AS red_fighter_draw,
		red_stats.win_by_ko AS red_fighter_ko_wins, red_stats.win_by_sub AS red_fighter_sub_wins, red_stats.win_by_dec AS red_fighter_dec_wins,
		f.fighter_blue_id, blue.name AS blue_fighter_name, blue.nickname AS blue_fighter_nickname, blue.division AS blue_fighter_division, blue.fighter_url AS blue_fighter_url,
		blue.image_url AS blue_fighter_image, blue.wins AS blue_fighter_wins, blue.loses AS blue_fighter_loses, blue.draw AS blue_fighter_draw,
		blue_stats.win_by_ko AS blue_fighter_ko_wins, blue_stats.win_by_sub AS blue_fighter_sub_wins, blue_stats.win_by_dec AS blue_fighter_dec_wins
	FROM
		filtered_events e
	LEFT JOIN
		fb_fights f ON e.event_id = f.event_id
	LEFT JOIN
		fb_fighters red ON f.fighter_red_id = red.fighter_id
	LEFT JOIN
		fb_fighters blue ON f.fighter_blue_id = blue.fighter_id
	LEFT JOIN
		fb_fighter_stats red_stats ON f.fighter_red_id = red_stats.fighter_id
	LEFT JOIN
		fb_fighter_stats blue_stats ON f.fighter_blue_id = blue_stats.fighter_id`

	var events []*model.FullEventResponse

	rows, err := r.Pool.Query(ctx, q, limit)
	if err != nil {
		return nil, r.DebugLogSqlErr(q, err)
	}
	defer rows.Close()

	fightMap := make(map[int32][]model.FightResponse)
	eventMap := make(map[int32]*model.FullEventResponse)

	for rows.Next() {
		var event model.FullEventResponse
		var fight model.FightResponse
		var redFighter model.Fighter
		var redFighterReq model.FighterReq
		var redFighterStats model.FighterStats
		var blueFighter model.Fighter
		var blueFighterReq model.FighterReq
		var blueFighterStats model.FighterStats

		if err := rows.Scan(
			&event.EventId, &event.Name, &event.IsDone, &fight.FightId, &fight.IsDone,
			&fight.NotContest, &fight.CreatedAt, &fight.FightDate, &fight.Result,
			&redFighterReq.FighterId, &redFighter.Name, &redFighter.NickName, &redFighter.Division, &redFighter.FighterUrl,
			&redFighter.ImageUrl, &redFighter.Wins, &redFighter.Loses, &redFighter.Draw,
			&redFighterStats.WinByKO, &redFighterStats.WinBySub, &redFighterStats.WinByDec,
			&blueFighterReq.FighterId, &blueFighter.Name, &blueFighter.NickName, &blueFighter.Division, &blueFighter.FighterUrl,
			&blueFighter.ImageUrl, &blueFighter.Wins, &blueFighter.Loses, &blueFighter.Draw,
			&blueFighterStats.WinByKO, &blueFighterStats.WinBySub, &blueFighterStats.WinByDec,
		); err != nil {
			return nil, r.DebugLogSqlErr(q, err)
		}

		redFighter.Stats = redFighterStats
		redFighterReq.Fighter = &redFighter
		blueFighter.Stats = blueFighterStats
		blueFighterReq.Fighter = &blueFighter

		fight.FighterRed = redFighterReq
		fight.FighterBlue = blueFighterReq

		fights, found := fightMap[event.EventId]
		if !found {
			fights = make([]model.FightResponse, 0)
		}

		fights = append(fights, fight)
		fightMap[event.EventId] = fights

		_, found = eventMap[event.EventId]
		if !found {
			eventMap[event.EventId] = &model.FullEventResponse{
				EventId: event.EventId,
				Name:    event.Name,
				IsDone:  event.IsDone,
			}
		}
	}

	for eid, event := range eventMap {
		event.Fights = fightMap[eid]
		events = append(events, event)
	}

	return events, nil
}

func (r *CommonRepo) GetEventId(ctx context.Context, tx pgx.Tx, fightId int32) (int32, error) {
	q := "SELECT event_id FROM fb_fights WHERE fight_id = $1"

	var eventId int32
	if tx != nil {
		if err := tx.QueryRow(ctx, q, fightId).Scan(&eventId); err != nil {
			return -1, r.DebugLogSqlErr(q, err)
		}
	} else {
		if err := r.Pool.QueryRow(ctx, q, fightId).Scan(&eventId); err != nil {
			return -1, r.DebugLogSqlErr(q, err)
		}
	}

	return eventId, nil
}

func (r *CommonRepo) GetUndoneFights(ctx context.Context, tx pgx.Tx, eventId int32) (int, error) {
	q := "SELECT COUNT(*) FROM fb_fights WHERE event_id = $1 AND is_done = false"
	var count int
	err := tx.QueryRow(ctx, q, eventId).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *CommonRepo) SetEventDone(ctx context.Context, tx pgx.Tx, eventId int32) error {
	q := "UPDATE fb_events SET is_done = true WHERE event_id = $1"

	_, err := tx.Exec(ctx, q, eventId)
	if err != nil {
		return err
	}

	return nil
}
