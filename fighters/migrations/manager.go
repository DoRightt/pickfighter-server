package migrations

import (
	"pickfighter.com/pkg/pgxs"
	"go.uber.org/zap"
)

const migrationTableCreate = `CREATE table(
	
	)`

type DBMigration struct {
	*pgxs.Repo
	logger *zap.SugaredLogger
	runs   bool
}

func (dbm *DBMigration) Shutdown() {
	if dbm.Repo != nil {
		dbm.Repo.GracefulShutdown()
	}
}
