package event

import (
	"context"

	internalErr "fightbettr.com/events/pkg/errors"
	eventmodel "fightbettr.com/events/pkg/model"
	logs "fightbettr.com/pkg/logger"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
)

// handleEventCreation creates a new event along with associated fights in a transaction.
// It takes the provided context, a transaction object, and the request containing event details.
// If successful, it returns a response with the newly created event information.
// If any error occurs during the process, it rolls back the transaction and returns an appropriate API error.
func (c *Controller) handleEventCreation(ctx context.Context, tx pgx.Tx, req *eventmodel.EventRequest) (*eventmodel.Event, error) {
	eventId, err := c.repo.TxCreateEvent(ctx, tx, req)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}

		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				intErr := internalErr.NewDefault(internalErr.TxNotUnique, 114)
				return nil, intErr
			}
		} else {
			intErr := internalErr.NewDefault(internalErr.TxUnknown, 115)
			logs.Errorf("Failed to create event during registration transaction: %s", err)
			return nil, intErr
		}
	}

	event := eventmodel.Event{
		EventId: eventId,
		Name:    req.Name,
		Fights:  req.Fights,
	}

	for _, f := range req.Fights {
		fight := eventmodel.Fight{
			EventId:       eventId,
			FighterRedId:  f.FighterRedId,
			FighterBlueId: f.FighterBlueId,
			IsDone:        false,
			IsCanceled:    false,
		}

		if err := c.repo.TxCreateEventFight(ctx, tx, fight); err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				logs.Errorf("Unable to rollback transaction: %s", txErr)
			}

			if pgErr, ok := err.(*pgconn.PgError); ok {
				if pgErr.Code == pgerrcode.UniqueViolation {
					intErr := internalErr.NewDefault(internalErr.TxNotUnique, 116)
					return nil, intErr
				}
			} else {
				intErr := internalErr.NewDefault(internalErr.TxUnknown, 117)
				logs.Errorf("Failed to create fight during registration transaction: %s", err)
				return nil, intErr
			}
		}
	}

	return &event, err
}

// checkEventIsDone checks if all fights are done. If so, sets event as done.
// It takes the fight ID as input and finds the corresponding event in which it is listed.
func (c *Controller) checkEventIsDone(ctx context.Context, tx pgx.Tx, fightId int32) error {
	eventId, err := c.repo.GetEventId(ctx, tx, fightId)
	if err != nil {
		return err
	}

	count, err := c.repo.GetUndoneFightsCount(ctx, tx, eventId)
	if err != nil {
		return err
	}

	if count == 0 {
		err = c.repo.SetEventDone(ctx, tx, eventId)
		if err != nil {
			return err
		}
	}

	return nil
}
