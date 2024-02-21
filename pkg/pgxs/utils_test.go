package pgxs

import (
	"errors"
	"projects/fb-server/pkg/logger"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestDebugLogSqlErr(t *testing.T) {
	lg := logger.NewSugared()

	tests := []struct {
		Name        string
		Query       string
		Error       error
		ExpectedLog string
	}{
		{
			Name:        "Error with PgError Code 23505",
			Query:       "SELECT * FROM table",
			Error:       &pgconn.PgError{Code: "23505", Message: "Some message"},
			ExpectedLog: "",
		},
		{
			Name:        "Error with PgError Code not 23505",
			Query:       "SELECT * FROM table",
			Error:       &pgconn.PgError{Code: "42P01", Message: "Table does not exist"},
			ExpectedLog: "query: \nSELECT * FROM table",
		},
		{
			Name:        "Error with non-PgError",
			Query:       "SELECT * FROM table",
			Error:       errors.New("Some generic error"),
			ExpectedLog: "query: \nSELECT * FROM table",
		},
		{
			Name:        "Error with ErrNoRows",
			Query:       "SELECT * FROM table",
			Error:       pgx.ErrNoRows,
			ExpectedLog: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			repo := &Repo{
				Logger: lg,
			}

			err := repo.DebugLogSqlErr(tc.Query, tc.Error)

			assert.Equal(t, tc.Error, err, "Unexpected error")
		})
	}
}

func TestSanitizeString(t *testing.T) {
	db := &Repo{}
	tests := []struct {
		Input    string
		Expected string
	}{
		{"Hello World", "Hello World"},
		{"Don't quote me", "Dont quote me"},
		{"100%", "100"},
		{"", ""},
		{"Let's test% this", "Lets test this"},
	}

	for _, tc := range tests {
		t.Run(tc.Input, func(t *testing.T) {
			result := db.SanitizeString(tc.Input)
			assert.Equal(t, tc.Expected, result, "Unexpected result for SanitizeString")
		})
	}
}

func TestQuoteString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "Hello World"},
		{"Don't quote me", "Dont quote me"},
		{"100%", "100"},
		{"", ""},
		{"Let's test% this", "Lets test this"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := QuoteString(tc.input)
			assert.Equal(t, tc.expected, result, "For input '%s', expected '%s' but got '%s'", tc.input, tc.expected, result)
		})
	}
}
