package event

import (
	"context"

	internalErr "github.com/DoRightt/pickfighter-server/events/pkg/errors"
	"github.com/DoRightt/pickfighter-server/events/pkg/model"
	logs "github.com/DoRightt/pickfighter-server/pkg/logger"
	"github.com/jackc/pgx/v5"
)

func (c *Controller) SetFightResult(ctx context.Context, req *model.FightResultRequest) (int32, error) {
	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		logs.Errorf("Unable to begin transaction: %s", err)
		intErr := internalErr.NewDefault(internalErr.Tx, 118)
		return 0, intErr
	}

	err = c.repo.CreateFightResult(ctx, tx, req)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		intErr := internalErr.New(internalErr.EventsFightResult, err, 904)
		return 0, intErr
	}

	err = c.repo.SetFightIsDone(ctx, tx, int(req.FightId))
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		intErr := internalErr.New(internalErr.EventsFightResult, err, 905)
		return 0, intErr
	}

	err = c.checkEventIsDone(ctx, tx, req.FightId)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		intErr := internalErr.New(internalErr.EventIsDone, err, 906)
		return 0, intErr
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		logs.Errorf("Unable to commit transaction: %s", txErr)
		intErr := internalErr.New(internalErr.TxCommit, txErr, 119)
		return 0, intErr
	}

	return req.FightId, nil
}
