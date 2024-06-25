package event

import (
	"context"

	internalErr "fightbettr.com/events/pkg/errors"
	eventmodel "fightbettr.com/events/pkg/model"
	"github.com/jackc/pgx/v5"
)

func (c *Controller) CreateBet(ctx context.Context, req *eventmodel.Bet) (int32, error) {
	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		c.Logger.Errorf("Unable to begin transaction: %s", err)
		cErr := internalErr.New(internalErr.Tx, err, 112)
		return 0, cErr
	}

	betId, err := c.repo.TxCreateBet(ctx, tx, req)
	if err != nil {
		c.Logger.Errorf("Error while user credentials creation: %s", err)
		return 0, err
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		c.Logger.Errorf("Unable to commit transaction: %s", err)
		cErr := internalErr.New(internalErr.TxCommit, err, 113)
		return 0, cErr
	}

	return betId, nil
}

func (c *Controller) GetBets(ctx context.Context, userId int32) (*eventmodel.BetsResponse, error) {
	count, err := c.repo.SearchBetsCount(ctx, userId)
	if err != nil {
		c.Logger.Errorf("Failed to get bets count: %s", err)
		intErr := internalErr.NewDefault(internalErr.BetsCount, 1201)

		return nil, intErr
	}

	if count == 0 {
		intErr := internalErr.NewDefault(internalErr.BetsNoRows, 1202)
		return nil, intErr
	}

	bets, err := c.repo.SearchBets(ctx, userId)
	if err != nil {
		c.Logger.Errorf("Failed to find bets: %s", err)
		intErr := internalErr.NewDefault(internalErr.Bets, 1203)
		return nil, intErr
	}

	return &eventmodel.BetsResponse{Bets: bets, Count: count}, nil
}
