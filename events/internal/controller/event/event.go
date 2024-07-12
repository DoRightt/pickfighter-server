package event

import (
	"context"

	internalErr "fightbettr.com/events/pkg/errors"
	"fightbettr.com/events/pkg/model"
	logs "fightbettr.com/pkg/logger"
	"github.com/jackc/pgx/v5"
)

func (c *Controller) CreateEvent(ctx context.Context, req *model.EventRequest) (int32, error) {
	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		logs.Errorf("Unable to begin transaction: %s", err)
		cErr := internalErr.New(internalErr.Tx, err, 112)
		return 0, cErr
	}

	event, err := c.handleEventCreation(ctx, tx, req)
	if err != nil {
		logs.Errorf("Error while user credentials creation: %s", err)
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		return 0, err
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		logs.Errorf("Unable to commit transaction: %s", err)
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		cErr := internalErr.New(internalErr.TxCommit, err, 113)
		return 0, cErr
	}

	return event.EventId, nil
}

func (c *Controller) GetEvents(ctx context.Context) (*model.EventsResponse, error) {
	count, err := c.repo.SearchEventsCount(ctx)
	if err != nil {
		logs.Errorf("Failed to get events count: %s", err)
		intErr := internalErr.NewDefault(internalErr.EventsCount, 901)

		return nil, intErr
	}
	if count == 0 {
		intErr := internalErr.NewDefault(internalErr.EventsNoRows, 902)
		return nil, intErr
	}

	events, err := c.repo.SearchEvents(ctx)
	if err != nil {
		logs.Errorf("Failed to find events: %s", err)
		intErr := internalErr.NewDefault(internalErr.Events, 903)
		return nil, intErr
	}
	return &model.EventsResponse{Count: count, Events: events}, nil
}
