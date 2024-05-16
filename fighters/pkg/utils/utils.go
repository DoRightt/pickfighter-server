package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"fightbettr.com/fighters/internal/repository/psql"
	lg "fightbettr.com/fighters/pkg/logger"
	"fightbettr.com/fighters/pkg/model"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
)

// ReadFighterData reads fighter data from a JSON file and returns a slice of model.Fighter.
// The file path is set to "../scraper/collection/fighters.json".
func ReadFighterData() ([]model.Fighter, error) {
	filePath := "../scraper/collection/fighters.json"

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("DEBUG", "ERROR FILE", err)
		return nil, err
	}

	var fightersData struct {
		Fighters []model.Fighter `json:"Fighters"`
	}

	if err := json.Unmarshal(jsonData, &fightersData); err != nil {
		return nil, err
	}

	return fightersData.Fighters, nil
}

// WriteFighterData writes fighter data to a PostgreSQL database using the provided context,
// and a slice of model.Fighter. It connects to the database using the configuration
// from ViperPostgres and performs create or update operations for each fighter.
func WriteFighterData(ctx context.Context, data []model.Fighter) error {
	l := lg.GetSugared()

	rep, err := psql.New(ctx, l)
	if err != nil {
		l.Errorf("Unable to start postgresql connection: %s", err)
	}
	defer rep.PoolClose()

	counter := 1

	for _, fighter := range data {
		fighterId, err := rep.FindFighter(ctx, fighter)
		if err != nil {
			if err == pgx.ErrNoRows {
				if err := createFighter(ctx, rep, fighter); err != nil {
					l.Errorf("Error while fighter transaction: %s", err)
					return err
				}

				fmt.Printf("[Operation №%d] Created: %s\n", counter, fighter.Name)
			} else {
				l.Errorf("Failed to find fighter: %s", err)
				return err
			}
		} else {
			fighter.FighterId = fighterId

			if err := updateFighter(ctx, rep, fighter); err != nil {
				l.Errorf("Error while fighter transaction: %s", err)
				return err
			}

			fmt.Printf("[Operation №%d] Updated: %s\n", counter, fighter.Name)
		}
		counter += 1
	}

	return nil
}

// DeleteFighterData deletes all records from the fb_fighters and fb_fighter_stats tables.
func DeleteFighterData(ctx context.Context) {
	l := lg.GetSugared()
	rep, err := psql.New(ctx, l)
	if err != nil {
		l.Errorf("Unable to start postgresql connection: %s", err)
	}

	fightersTableNames := []string{"fb_fighters", "fb_fighter_stats"}
	handledTableNames := []string{}

	for _, name := range fightersTableNames {
		err = rep.DeleteRecords(ctx, name)
		if err != nil {
			l.Fatalf("Error deleting records: %s", err)
		}

		handledTableNames = append(handledTableNames, name)
	}

	for _, name := range handledTableNames {
		fmt.Printf("All records from table '%s' deleted successfully\n", name)
	}
}

// createNewFighterTx performs a transaction to create a new fighter in the database.
// It takes a context, a fighter repository, and a model.Fighter as parameters.
// If the transaction fails, it logs the error and returns an appropriate ApiError.
func createFighter(ctx context.Context, rep *psql.Repository, fighter model.Fighter) error {
	l := lg.GetSugared()
	tx, err := rep.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		l.Errorf("Unable to begin transaction: %s", err)
	}
	defer tx.Rollback(ctx)

	fighterId, err := rep.CreateNewFighter(ctx, tx, fighter)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			l.Errorf("Unable to rollback transaction: %s", txErr)
		}
		if err.(*pgconn.PgError).Code == pgerrcode.UniqueViolation {
			// intErr := internalErr.New(internalErr.TxNotUnique)
			// return httplib.NewApiErrFromInternalErr(intErr)
		} else {
			// intErr := internalErr.New(internalErr.TxUnknown)
			l.Errorf("Failed to create fighter during registration transaction: %s", err)
			// return httplib.NewApiErrFromInternalErr(intErr, http.StatusInternalServerError)
		}
	}

	fighter.Stats.FighterId = fighterId
	err = rep.CreateNewFighterStats(ctx, tx, fighter.Stats)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			l.Errorf("Unable to rollback transaction: %s", txErr)
		}
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		l.Errorf("Unable to commit transaction: %s", txErr)
		return txErr
	}

	return nil
}

// updateFighterTx performs a transaction to update an existing fighter in the database.
// It takes a context, a fighter repository, and a model.Fighter as parameters.
// If the transaction fails, it logs the error and returns an appropriate ApiError.
func updateFighter(ctx context.Context, rep *psql.Repository, fighter model.Fighter) error {
	l := lg.GetSugared()
	tx, err := rep.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		l.Errorf("Unable to begin transaction: %s", err)
	}
	defer tx.Rollback(ctx)

	updatedId, err := rep.UpdateFighter(ctx, tx, fighter)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			l.Errorf("Unable to rollback transaction: %s", txErr)
		}
		if err.(*pgconn.PgError).Code == pgerrcode.UniqueViolation {
			// intErr := internalErr.New(internalErr.TxNotUnique)
			// return httplib.NewApiErrFromInternalErr(intErr)
		} else {
			// intErr := internalErr.New(internalErr.TxUnknown)
			l.Errorf("Failed to update fighter during registration transaction: %s", err)
			// return httplib.NewApiErrFromInternalErr(intErr, http.StatusInternalServerError)
		}
	}

	fighter.Stats.FighterId = updatedId

	if err := rep.UpdateFighterStats(ctx, tx, fighter.Stats); err != nil {
		l.Errorf("Error while updating stats: %s", err)
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		l.Errorf("Unable to commit transaction: %s", txErr)
		return txErr
	}

	return nil
}
