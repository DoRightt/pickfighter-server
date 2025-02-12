package migrations

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/DoRightt/pickfighter-server/events/internal/repository/psql"
	"github.com/jackc/pgx/v5"
)

func InitEventsSchema(ctx context.Context, r *psql.Repository) error {
	tx, err := r.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	queries := [][]string{
		{
			// Queries to create Schema
			`CREATE SCHEMA IF NOT EXISTS events;`,
		},
		{
			// Queries to create tables
			`CREATE TABLE IF NOT EXISTS events.bets (
				bet_id integer NOT NULL,
				user_id integer,
				fight_id integer,
				created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
				bet integer NOT NULL
			);`,
			`CREATE TABLE IF NOT EXISTS events.events (
				event_id integer NOT NULL,
				name character varying(255) NOT NULL,
				is_done boolean DEFAULT false
			);`,
			`CREATE TABLE IF NOT EXISTS events.fight_results (
				result_id integer NOT NULL,
				fight_id integer,
				winner_id integer,
				not_contest boolean DEFAULT false,
				is_draw boolean DEFAULT false,
				round smallint,
				method VARCHAR(50),
				created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
			);`,
			`CREATE TABLE IF NOT EXISTS events.fights (
				fight_id integer NOT NULL,
				fighter_red_id integer NOT NULL,
				fighter_blue_id integer NOT NULL,
				is_done boolean DEFAULT false,
				created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
				fight_date timestamp without time zone,
				is_canceled boolean DEFAULT false,
				event_id integer
			);`,
		},
		{
			// Queries to create sequences
			`CREATE SEQUENCE IF NOT EXISTS events.pf_bets_bet_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;`,
			`CREATE SEQUENCE IF NOT EXISTS events.pf_events_event_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;`,
			`CREATE SEQUENCE IF NOT EXISTS events.pf_fight_results_result_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;`,
			`CREATE SEQUENCE IF NOT EXISTS events.pf_fights_fight_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;`,
		},
		{
			// Queries to set sequences owners
			`ALTER SEQUENCE events.pf_bets_bet_id_seq OWNED BY events.bets.bet_id;`,
			`ALTER SEQUENCE events.pf_events_event_id_seq OWNED BY events.events.event_id;`,
			`ALTER SEQUENCE events.pf_fight_results_result_id_seq OWNED BY events.fight_results.result_id;`,
			`ALTER SEQUENCE events.pf_fights_fight_id_seq OWNED BY events.fights.fight_id;`,
		},
		{
			// Queries to set default values
			`ALTER TABLE ONLY events.bets ALTER COLUMN bet_id SET DEFAULT nextval('events.pf_bets_bet_id_seq'::regclass);`,
			`ALTER TABLE ONLY events.events ALTER COLUMN event_id SET DEFAULT nextval('events.pf_events_event_id_seq'::regclass);`,
			`ALTER TABLE ONLY events.fight_results ALTER COLUMN result_id SET DEFAULT nextval('events.pf_fight_results_result_id_seq'::regclass);`,
			`ALTER TABLE ONLY events.fights ALTER COLUMN fight_id SET DEFAULT nextval('events.pf_fights_fight_id_seq'::regclass);`,
		},
		{
			// Queries to set constraints
			`ALTER TABLE ONLY events.bets ADD CONSTRAINT pf_bets_pk PRIMARY KEY (bet_id);`,
			`ALTER TABLE ONLY events.events ADD CONSTRAINT pf_events_pkey PRIMARY KEY (event_id);`,
			`ALTER TABLE ONLY events.fight_results ADD CONSTRAINT pf_fight_results_pk PRIMARY KEY (result_id);`,
			`ALTER TABLE ONLY events.fights ADD CONSTRAINT pf_fights_pk PRIMARY KEY (fight_id);`,
			`ALTER TABLE ONLY events.bets ADD CONSTRAINT pf_bets_event_id_fkey FOREIGN KEY (fight_id) REFERENCES events.fights(fight_id);`,
			`ALTER TABLE ONLY events.bets ADD CONSTRAINT pf_bets_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth.users(user_id);`,
			`ALTER TABLE ONLY events.fight_results ADD CONSTRAINT pf_fight_results_fb_fighters_fighter_id_fk FOREIGN KEY (winner_id) REFERENCES fighters.fighters(fighter_id);`,
			`ALTER TABLE ONLY events.fight_results ADD CONSTRAINT pf_fight_results_fb_fights_fight_id_fk FOREIGN KEY (fight_id) REFERENCES events.fights(fight_id);`,
			`ALTER TABLE ONLY events.fights ADD CONSTRAINT pf_fights_fb_events_event_id_fk FOREIGN KEY (event_id) REFERENCES events.events(event_id);`,
			`ALTER TABLE ONLY events.fights ADD CONSTRAINT pf_fights_fighter_blue_id_fkey FOREIGN KEY (fighter_blue_id) REFERENCES fighters.fighters(fighter_id);`,
			`ALTER TABLE ONLY events.fights ADD CONSTRAINT pf_fights_fighter_red_id_fkey FOREIGN KEY (fighter_red_id) REFERENCES fighters.fighters(fighter_id);`,
		},
	}

	for _, block := range queries {
		for _, query := range block {

			if _, err = tx.Exec(ctx, query); err != nil {

				if strings.Contains(err.Error(), "already exists") {
					log.Println("Object already exists, skipping...")
				} else {
					if txErr := tx.Rollback(ctx); txErr != nil {
						fmt.Printf("Unable to rollback transaction: %s", txErr)
					}
					return fmt.Errorf("failed to execute migration query: %w", err)
				}

			}

		}
	}

	if err = tx.Commit(ctx); err != nil {
		fmt.Printf("Unable to commit transaction: %s", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Println("Migrations applied successfully")

	return nil
}
