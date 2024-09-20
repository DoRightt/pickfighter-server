package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"pickfighter.com/events/pkg/model"
	eventmodel "pickfighter.com/events/pkg/model"
	"pickfighter.com/events/pkg/version"
	"pickfighter.com/pkg/pgxs"
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

// HealthCheck returns the current health status of the application.
// It includes information such as the app version, start time, uptime,
// and a message indicating the application's health.
func (c *Controller) HealthCheck() *model.HealthStatus {
	return &model.HealthStatus{
		AppDevVersion: version.DevVersion,
		AppName:       version.Name,
		Timestamp:     time.Now().Format(time.RFC1123),
		AppRunDate:    version.RunDate,
		AppTimeAlive:  time.Now().Unix() - version.RunDate,
		Healthy:       true,
		Message:       fmt.Sprintf("[%s]: I'm fine, thanks!", viper.GetString("app.name")),
	}
}
