package scraper

import (
	"context"
	"encoding/json"
	"os"
	"projects/fb-server/pkg/cfg"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/pgxs"
	fighterRepo "projects/fb-server/repo/fighters"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func ReadFighterData() ([]model.Fighter, error) {
	filePath := "data/fighters.json"

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
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

func WriteFighterData(ctx context.Context, l *zap.SugaredLogger, data []model.Fighter) error {
	db, err := pgxs.NewPool(ctx, l, cfg.ViperPostgres())
	if err != nil {
		l.Fatalf("Unable to connect postgresql: %s", err)
	}
	defer db.Pool.Close()

	rep := fighterRepo.New(db)

	tx, err := rep.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		l.Errorf("Unable to begin transaction: %s", err)
	}
	defer tx.Rollback(ctx)

	err = rep.InsertFightersData(ctx, tx, data)

	if txErr := tx.Commit(ctx); txErr != nil {
		l.Errorf("Unable to commit transaction: %s", txErr)
		return txErr
	}

	return nil
}
