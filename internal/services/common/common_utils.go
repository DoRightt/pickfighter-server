package common

import (
	"context"
	"net/http"
	internalErr "projects/fb-server/pkg/errors"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// CreateEvent creates a new event along with associated fights in a transaction.
// It takes the provided context, a transaction object, and the request containing event details.
// If successful, it returns a response with the newly created event information.
// If any error occurs during the process, it rolls back the transaction and returns an appropriate API error.
func (s *service) CreateEvent(ctx context.Context, tx pgx.Tx, req *model.EventsRequest) (*model.EventResponse, error) {
	eventId, err := s.Repo.TxCreateEvent(ctx, tx, req)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			s.Logger.Errorf("Unable to rollback transaction: %s", txErr)
		}

		if err.(*pgconn.PgError).Code == pgerrcode.UniqueViolation {
			intErr := internalErr.New(internalErr.TxNotUnique)
			return nil, httplib.NewApiErrFromInternalErr(intErr)
		} else {
			intErr := internalErr.New(internalErr.TxUnknown)
			s.Logger.Errorf("Failed to create event during registration transaction: %s", err)
			return nil, httplib.NewApiErrFromInternalErr(intErr, http.StatusInternalServerError)
		}
	}

	event := model.EventResponse{
		EventId: eventId,
		Name:    req.Name,
		Fights:  req.Fights,
	}

	for _, f := range req.Fights {
		fight := model.Fight{
			EventId:       eventId,
			FighterRedId:  f.FighterRedId,
			FighterBlueId: f.FighterBlueId,
			IsDone:        false,
			IsCanceled:    false,
		}

		if err := s.Repo.TxCreateEventFight(ctx, tx, fight); err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				s.Logger.Errorf("Unable to rollback transaction: %s", txErr)
			}
			if err.(*pgconn.PgError).Code == pgerrcode.UniqueViolation {
				intErr := internalErr.New(internalErr.TxNotUnique)
				return nil, httplib.NewApiErrFromInternalErr(intErr)
			} else {
				intErr := internalErr.New(internalErr.TxUnknown)
				s.Logger.Errorf("Failed to create fight during registration transaction: %s", err)
				return nil, httplib.NewApiErrFromInternalErr(intErr, http.StatusInternalServerError)
			}
		}
	}

	return &event, err
}

// CheckEventIsDone checks if all fights are done. If so, sets event as done.
// It takes the fight ID as input and finds the corresponding event in which it is listed. 
func (s *service) CheckEventIsDone(ctx context.Context, tx pgx.Tx, fightId int32) error {
	eventId, err := s.Repo.GetEventId(ctx, tx, fightId)
	if err != nil {
		return err
	}

	count, err := s.Repo.GetUndoneFightsCount(ctx, tx, eventId)
	if err != nil {
		return err
	}

	if count == 0 {
		err = s.Repo.SetEventDone(ctx, tx, eventId)
		if err != nil {
			return err
		}
	}

	return nil
}

// capitalize returns a capitalized string.
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
