package event

import (
	"context"
	"errors"

	eventmodel "pickfighter.com/events/pkg/model"
	"pickfighter.com/pkg/pgxs"
	"github.com/jackc/pgx/v5"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type eventRepository interface {
	pgxs.PickfighterRepo

	TxCreateEvent(ctx context.Context, tx pgx.Tx, e *eventmodel.EventRequest) (int32, error)
	TxCreateEventFight(ctx context.Context, tx pgx.Tx, f eventmodel.Fight) error
	SearchEventsCount(ctx context.Context) (int32, error)
	SearchEvents(ctx context.Context) ([]*eventmodel.Event, error)
	TxCreateBet(ctx context.Context, tx pgx.Tx, req *eventmodel.Bet) (int32, error)
	SearchBetsCount(ctx context.Context, userId int32) (int32, error)
	SearchBets(ctx context.Context, userId int32) ([]*eventmodel.Bet, error)
	SetFightResult(ctx context.Context, tx pgx.Tx, fr *eventmodel.FightResultRequest) error
	GetEventId(ctx context.Context, tx pgx.Tx, fightId int32) (int32, error)
	GetUndoneFightsCount(ctx context.Context, tx pgx.Tx, eventId int32) (int, error)
	SetEventDone(ctx context.Context, tx pgx.Tx, eventId int32) error
}

// Controller defines a metadata service controller.
type Controller struct {
	repo eventRepository
}

// New creates a Event service controller.
func New(repo eventRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}
