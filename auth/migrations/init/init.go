package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/DoRightt/pickfighter-server/auth/internal/repository/psql"
)

func InitAuthSchema(ctx context.Context, r *psql.Repository) error {
	tx, err := r.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	queries := []string{
		`CREATE SCHEMA IF NOT EXISTS auth;`,

		`CREATE TABLE IF NOT EXISTS auth.users (
			user_id SERIAL PRIMARY KEY,
			name VARCHAR NOT NULL,
			rank VARCHAR(60),
			claim VARCHAR(60) DEFAULT 'USER',
			flags INTEGER NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL DEFAULT (EXTRACT(epoch FROM now()))::INTEGER,
			updated_at INTEGER
		);`,

		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name ON auth.users(name);`,

		`CREATE TABLE IF NOT EXISTS auth.user_credentials (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			email VARCHAR NOT NULL UNIQUE,
			password_hash VARCHAR NOT NULL,
			salt VARCHAR NOT NULL,
			token VARCHAR UNIQUE,
			token_type VARCHAR,
			active BOOLEAN DEFAULT false,
			token_expire INTEGER,
			created_at INTEGER NOT NULL DEFAULT (EXTRACT(epoch FROM now()))::INTEGER,
			updated_at INTEGER,
			CONSTRAINT fk_user_credentials_user_id FOREIGN KEY (user_id) 
				REFERENCES auth.users(user_id) ON DELETE CASCADE
		);`,

		`CREATE UNIQUE INDEX IF NOT EXISTS idx_user_credentials_email ON auth.user_credentials(email);`,

		`CREATE INDEX IF NOT EXISTS idx_user_credentials_user_id ON auth.user_credentials(user_id);`,
	}

	for _, query := range queries {
		if _, err = tx.Exec(ctx, query); err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				fmt.Printf("Unable to rollback transaction: %s", txErr)
			}
			return fmt.Errorf("failed to execute migration query: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		fmt.Printf("Unable to commit transaction: %s", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Println("Migrations applied successfully")

	return nil
}
