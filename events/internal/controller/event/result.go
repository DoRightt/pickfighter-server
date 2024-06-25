package event

import (
	"context"

	internalErr "fightbettr.com/events/pkg/errors"
	"fightbettr.com/events/pkg/model"
	"github.com/jackc/pgx/v5"
)

func (c *Controller) SetFightResult(ctx context.Context, req *model.FightResultRequest) (int32, error) {
	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		c.Logger.Errorf("Unable to begin transaction: %s", err)
		intErr := internalErr.NewDefault(internalErr.Tx, 118)
		return 0, intErr
	}

	err = c.repo.SetFightResult(ctx, tx, req)
	if err != nil {
		intErr := internalErr.New(internalErr.EventsFightResult, err, 904)
		return 0, intErr
	}

	err = c.checkEventIsDone(ctx, tx, req.FightId)
	if err != nil {
		intErr := internalErr.New(internalErr.EventIsDone, err, 905)
		return 0, intErr
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		c.Logger.Errorf("Unable to commit transaction: %s", txErr)
		intErr := internalErr.New(internalErr.TxCommit, txErr, 119)
		return 0, intErr
	}

	return req.FightId, nil
}
