package event

import (
	"context"

	internalErr "pickfighter.com/events/pkg/errors"
	eventmodel "pickfighter.com/events/pkg/model"
	logs "pickfighter.com/pkg/logger"
	"github.com/jackc/pgx/v5"
)

func (c *Controller) CreateBet(ctx context.Context, req *eventmodel.Bet) (int32, error) {
	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		logs.Errorf("Unable to begin transaction: %s", err)
		cErr := internalErr.New(internalErr.Tx, err, 112)
		return 0, cErr
	}

	betId, err := c.repo.TxCreateBet(ctx, tx, req)
	if err != nil {
		logs.Errorf("Error while user credentials creation: %s", err)
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		return 0, err
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		logs.Errorf("Unable to commit transaction: %s", err)
		cErr := internalErr.New(internalErr.TxCommit, err, 113)
		return 0, cErr
	}

	return betId, nil
}

func (c *Controller) GetBets(ctx context.Context, userId int32) (*eventmodel.BetsResponse, error) {
	count, err := c.repo.SearchBetsCount(ctx, userId)
	if err != nil {
		logs.Errorf("Failed to get bets count: %s", err)
		intErr := internalErr.NewDefault(internalErr.BetsCount, 1201)

		return nil, intErr
	}

	if count == 0 {
		intErr := internalErr.NewDefault(internalErr.BetsNoRows, 1202)
		return nil, intErr
	}

	bets, err := c.repo.SearchBets(ctx, userId)
	if err != nil {
		logs.Errorf("Failed to find bets: %s", err)
		intErr := internalErr.NewDefault(internalErr.Bets, 1203)
		return nil, intErr
	}

	return &eventmodel.BetsResponse{Bets: bets, Count: count}, nil
}
