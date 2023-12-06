package common

import (
	"context"
	"net/http"
	internalErr "projects/fb-server/errors"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

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

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
