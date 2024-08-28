package pgxs

import (
	"fmt"
	"strings"

	logs "pickfighter.com/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrEmptyConfig is returned when the PSQL Config is empty.
var ErrEmptyConfig = fmt.Errorf("pxgs: PSQL Config is required")

// DebugLogSqlErr logs debug information for SQL queries and errors.
// It takes the SQL query string 'q' and the error 'err' as parameters.
// If the error is a PostgreSQL "23505" violation, it sets the 'deuce' flag.
// It logs a debug message with the SQL query if the error is not of type pgx.ErrNoRows
// and not a "23505" violation.
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
		logs.Debugf("query: \n%s", q)
	}

	return err
}

// SanitizeString removes single quotes (') and percentage signs (%) from the given string.
func (db *Repo) SanitizeString(s string) string {
	return QuoteString(s)
}

// QuoteString removes single quotes (') and percentage signs (%) from the given string.
// It replaces single quotes and percentage signs with an empty string to prevent SQL injection issues.
// Returns the modified string.
func QuoteString(str string) string {
	str = strings.Replace(str, "'", "", -1)
	str = strings.Replace(str, "%", "", -1)
	return str
}
