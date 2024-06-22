package event

import (
	"context"
	"net/http"

	internalErr "fightbettr.com/events/pkg/errors"
	eventmodel "fightbettr.com/events/pkg/model"
	"fightbettr.com/pkg/httplib"
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
			c.Logger.Errorf("Unable to rollback transaction: %s", txErr)
		}

		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				intErr := internalErr.NewDefault(internalErr.TxNotUnique, 114)
				return nil, httplib.NewApiErrFromInternalErr(intErr)
			}
		} else {
			intErr := internalErr.NewDefault(internalErr.TxUnknown, 115)
			c.Logger.Errorf("Failed to create event during registration transaction: %s", err)
			return nil, httplib.NewApiErrFromInternalErr(intErr, http.StatusInternalServerError)
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
				c.Logger.Errorf("Unable to rollback transaction: %s", txErr)
			}

			if pgErr, ok := err.(*pgconn.PgError); ok {
				if pgErr.Code == pgerrcode.UniqueViolation {
					intErr := internalErr.NewDefault(internalErr.TxNotUnique, 116)
					return nil, httplib.NewApiErrFromInternalErr(intErr)
				}
			} else {
				intErr := internalErr.NewDefault(internalErr.TxUnknown, 117)
				c.Logger.Errorf("Failed to create fight during registration transaction: %s", err)
				return nil, httplib.NewApiErrFromInternalErr(intErr, http.StatusInternalServerError)
			}
		}
	}

	return &event, err
}
