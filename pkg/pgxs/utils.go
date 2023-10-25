package pgxs

import (
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

var ErrEmptyConfig = fmt.Errorf("pxgs: PSQL Config is required")

func (db *Repo) DebugLogSqlErr(q string, err error) error {
	var deuce bool
	fmt.Println("ERROR: ", err)
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		if pgErr.Code == "23505" {
			deuce = true
		}
	}

	if err != pgx.ErrNoRows && !deuce {
		db.Logger.Debugf("query: \n%s", q)
	}

	return err
}
