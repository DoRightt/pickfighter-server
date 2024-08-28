package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
	eventmodel "pickfighter.com/events/pkg/model"
)

// TxCreateEvent creates a new event in the 'pf_events' table and returns the event ID.
// It uses a transaction (tx) if provided, otherwise, it uses the repository's connection pool.
func (r *Repository) TxCreateEvent(ctx context.Context, tx pgx.Tx, e *eventmodel.EventRequest) (int32, error) {
	q := `INSERT INTO public.pf_events 
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
		if err := r.GetPool().QueryRow(ctx, q, args...).Scan(&eventId); err != nil {
			return 0, r.DebugLogSqlErr(q, err)
		}
	}

	return eventId, nil
}

// SearchEventsCount returns the count of events in the system, considering the limit constraint.
func (r *Repository) SearchEventsCount(ctx context.Context) (int32, error) {
	limit := 5

	q := ` WITH ranked_events AS (
		SELECT
			*,
			ROW_NUMBER() OVER (PARTITION BY is_done ORDER BY event_id) AS rn
		FROM public.pf_events
		WHERE is_done IN (false, true)
	),
	filtered_events AS (
		SELECT
			event_id,
			is_done
		FROM ranked_events
		WHERE (is_done = false AND rn <= $1) OR (is_done = true AND rn > (SELECT COUNT(*) FROM public.pf_events WHERE is_done = true) - $1)
	)
	SELECT COUNT(*) FROM filtered_events`

	var count int32
	if err := r.GetPool().QueryRow(ctx, q, limit).Scan(&count); err != nil {
		return 0, r.DebugLogSqlErr(q, err)
	}

	return count, nil
}

// SearchEvents retrieves a list of events along with associated fights and fighter information.
// It uses Common Table Expressions (CTE) to rank events and filter them based on whether they are done or not.
// The method takes a context, retrieves a limited number of events (specified by the 'limit' parameter),
// and returns a slice of FullEventResponse containing event details, associated fights, and fighter information.
// If there is an error during the database query, it returns nil and the encountered error.
func (r *Repository) SearchEvents(ctx context.Context) ([]*eventmodel.Event, error) {
	limit := 5

	q := `WITH ranked_events AS (
		SELECT
			*,
			ROW_NUMBER() OVER (PARTITION BY is_done ORDER BY event_id) AS rn
		FROM public.pf_events
		WHERE is_done IN (false, true)
	),
	filtered_events AS (
		SELECT
			event_id, name, is_done
		FROM ranked_events
		WHERE (is_done = false AND rn <= $1) OR (is_done = true AND rn > (SELECT COUNT(*) FROM public.pf_events WHERE is_done = true) - $1)
	)
	SELECT
		e.event_id, e.name, e.is_done AS is_event_done, 
		f.fight_id, f.is_done AS is_fight_done, f.not_contest, 
		f.created_at, f.fight_date, f.result,
		f.fighter_red_id, f.fighter_blue_id
	FROM
		filtered_events e
	LEFT JOIN
		pf_fights f ON e.event_id = f.event_id`

	var events []*eventmodel.Event

	rows, err := r.GetPool().Query(ctx, q, limit)
	if err != nil {
		return nil, r.DebugLogSqlErr(q, err)
	}
	defer rows.Close()

	fightMap := make(map[int32][]eventmodel.Fight)
	eventMap := make(map[int32]*eventmodel.Event)

	for rows.Next() {
		var event eventmodel.Event
		var fight eventmodel.Fight

		if err := rows.Scan(
			&event.EventId, &event.Name, &event.IsDone,
			&fight.FightId, &fight.IsDone, &fight.NotContest,
			&fight.CreatedAt, &fight.FightDate, &fight.Result,
			&fight.FighterRedId, &fight.FighterBlueId,
		); err != nil {
			return nil, r.DebugLogSqlErr(q, err)
		}

		fights, found := fightMap[event.EventId]
		if !found {
			fights = make([]eventmodel.Fight, 0)
		}

		fights = append(fights, fight)
		fightMap[event.EventId] = fights

		_, found = eventMap[event.EventId]
		if !found {
			eventMap[event.EventId] = &eventmodel.Event{
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

// GetEventId retrieves the event ID associated with a specific fight from the pf_fights table.
// It takes a transaction (tx), the fight ID, and returns the corresponding event ID.
// If the query is successful, it returns the event ID; otherwise, it returns -1 and the encountered error.
func (r *Repository) GetEventId(ctx context.Context, tx pgx.Tx, fightId int32) (int32, error) {
	q := "SELECT event_id FROM pf_fights WHERE fight_id = $1"

	var eventId int32
	if tx != nil {
		if err := tx.QueryRow(ctx, q, fightId).Scan(&eventId); err != nil {
			return -1, r.DebugLogSqlErr(q, err)
		}
	} else {
		if err := r.GetPool().QueryRow(ctx, q, fightId).Scan(&eventId); err != nil {
			return -1, r.DebugLogSqlErr(q, err)
		}
	}

	return eventId, nil
}

// GetUndoneFightsCount retrieves the count of undone fights for a specific event from the pf_fights table.
// It takes a transaction (tx), the event ID, and returns the number of fights that are not marked as done (is_done = false).
// If the query is successful, it returns the count, otherwise, it returns an error.
func (r *Repository) GetUndoneFightsCount(ctx context.Context, tx pgx.Tx, eventId int32) (int, error) {
	q := "SELECT COUNT(*) FROM pf_fights WHERE event_id = $1 AND is_done = false"
	var count int
	err := tx.QueryRow(ctx, q, eventId).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

// SetEventDone updates the 'is_done' field of an event in the pf_events table.
// It takes a transaction (tx) and the event ID as parameters and sets the 'is_done'
// column to true for the specified event. If the update is successful, it returns nil.
// In case of an error during the update, it returns the error details.
func (r *Repository) SetEventDone(ctx context.Context, tx pgx.Tx, eventId int32) error {
	q := "UPDATE pf_events SET is_done = true WHERE event_id = $1"

	_, err := tx.Exec(ctx, q, eventId)
	if err != nil {
		return err
	}

	return nil
}
