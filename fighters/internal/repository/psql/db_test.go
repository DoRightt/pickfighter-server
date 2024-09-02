package psql

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"pickfighter.com/fighters/pkg/cfg"
	"pickfighter.com/fighters/pkg/model"
)

var testFighter = &model.Fighter{
	FighterId:      57913,
	Name:           "Fabio Agu",
	Status:         "Active",
	FighterUrl:     "https://www.ufc.com/athlete/fabio-agu",
	DebutTimestamp: 1715817600,
}

func TestNew(t *testing.T) {
	repo, _ := New(context.Background(), cfg.ViperPostgres())

	if repo == nil {
		t.Errorf("expected repo not to be nil, got %v", repo)
	}
}

func TestPerformFightersQuery(t *testing.T) {
	tests := []struct {
		name     string
		req      *model.FightersRequest
		expected []string
	}{
		{
			name:     "nil request",
			req:      nil,
			expected: []string{},
		},
		{
			name: "status only",
			req: &model.FightersRequest{
				Status: "active",
			},
			expected: []string{
				`f.status = 'active'`,
			},
		},
		{
			name: "fighters IDs only",
			req: &model.FightersRequest{
				FightersIds: []int32{1, 2, 3},
			},
			expected: []string{
				`f.fighter_id IN (1, 2, 3)`,
			},
		},
		{
			name: "status and fighters IDs",
			req: &model.FightersRequest{
				Status:      "inactive",
				FightersIds: []int32{4, 5},
			},
			expected: []string{
				`f.status = 'inactive'`,
				`f.fighter_id IN (4, 5)`,
			},
		},
		{
			name: "empty status and empty fighters IDs",
			req: &model.FightersRequest{
				Status:      "",
				FightersIds: nil,
			},
			expected: []string{},
		},
	}

	repo := &Repository{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := repo.performFightersQuery(tc.req)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}

func TestSearchFightersCount(t *testing.T) {
	initTestConfig()
	defer viper.Reset()

	ctx := context.Background()
	config := cfg.ViperTestPostgres()

	repo, err := New(ctx, config)
	assert.NoError(t, err)
	defer repo.GracefulShutdown()

	count, err := repo.SearchFightersCount(ctx, &model.FightersRequest{Status: "Active"})

	assert.Greater(t, int(count), 0)
	assert.NoError(t, err)
}

func TestSearchFighters(t *testing.T) {
	initTestConfig()
	defer viper.Reset()

	ctx := context.Background()
	config := cfg.ViperTestPostgres()

	repo, err := New(ctx, config)
	assert.NoError(t, err)
	defer repo.GracefulShutdown()

	fighters, err := repo.SearchFighters(ctx, &model.FightersRequest{Status: "Active", FightersIds: []int32{testFighter.FighterId}})
	assert.Equal(t, fighters[0].FighterId, testFighter.FighterId)
	assert.NoError(t, err)
}

func TestFindFighter(t *testing.T) {
	initTestConfig()
	defer viper.Reset()

	ctx := context.Background()
	config := cfg.ViperTestPostgres()

	repo, err := New(ctx, config)
	assert.NoError(t, err)
	defer repo.GracefulShutdown()

	fighterIO, err := repo.FindFighter(ctx, *testFighter)
	assert.Equal(t, fighterIO, testFighter.FighterId)
	assert.NoError(t, err)
}

func TestCreateNewFighter(t *testing.T) {
	initTestConfig()
	defer viper.Reset()

	ctx := context.Background()
	config := cfg.ViperTestPostgres()

	repo, err := New(ctx, config)
	assert.NoError(t, err)
	defer repo.GracefulShutdown()

	tests := []struct {
		name          string
		fighter       model.Fighter
		expectedError *pgconn.PgError
	}{
		{
			name: "Success",
			fighter: model.Fighter{
				Name:           "Noob Saibot",
				NickName:       "",
				Division:       1,
				Status:         "Active",
				Hometown:       "",
				Height:         100,
				Weight:         100,
				OctagonDebut:   time.Now().Format(time.RFC1123),
				DebutTimestamp: int(time.Now().Unix()),
				Reach:          120,
				LegReach:       120,
				Wins:           1,
				Loses:          1,
				Draw:           1,
				FighterUrl:     "http://" + time.Now().Format(time.RFC1123),
				ImageUrl:       "",
			},
			expectedError: nil,
		},
		{
			name: "Error duplicate fighter url",
			fighter: model.Fighter{
				Name:           "Tobias Boon",
				NickName:       "",
				Division:       1,
				Status:         "Active",
				Hometown:       "",
				Height:         100,
				Weight:         100,
				OctagonDebut:   time.Now().Format(time.RFC1123),
				DebutTimestamp: int(time.Now().Unix()),
				Reach:          120,
				LegReach:       120,
				Wins:           1,
				Loses:          1,
				Draw:           1,
				FighterUrl:     "https://www.ufc.com/athlete/tank-abbott",
				ImageUrl:       "",
			},
			expectedError: &pgconn.PgError{
				Severity: "ERROR",
				Code:     "23505",
				Message:  "duplicate key value violates unique constraint \"pf_fighters_fighter_url_uindex\"",
				Detail:   "Key (fighter_url)=(https://www.ufc.com/athlete/tank-abbott) already exists.",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fighterIO, err := repo.CreateNewFighter(ctx, nil, tc.fighter)

			if err != nil {
				var pgErr *pgconn.PgError
				if assert.ErrorAs(t, err, &pgErr) {
					assert.Equal(t, tc.expectedError.Code, pgErr.Code)
					assert.Equal(t, tc.expectedError.Message, pgErr.Message)
				} else {
					t.Fatalf("expected error of type *pgconn.PgError but got %T", err)
				}
			} else {
				assert.Greater(t, fighterIO, int32(0))
			}
		})
	}
}

func TestCreateNewFighterStats(t *testing.T) {
	initTestConfig()
	defer viper.Reset()

	ctx := context.Background()
	config := cfg.ViperTestPostgres()

	repo, err := New(ctx, config)
	assert.NoError(t, err)
	defer repo.GracefulShutdown()

	err = repo.CreateNewFighterStats(ctx, nil, model.FighterStats{FighterId: testFighter.FighterId})
	assert.NoError(t, err)

	tx, err := repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	assert.NoError(t, err)

	err = repo.CreateNewFighterStats(ctx, tx, model.FighterStats{FighterId: testFighter.FighterId})
	assert.NoError(t, err)

	err = tx.Commit(ctx)
	assert.NoError(t, err)
}

func TestUpdateFighter(t *testing.T) {
	initTestConfig()
	defer viper.Reset()

	ctx := context.Background()
	config := cfg.ViperTestPostgres()

	repo, err := New(ctx, config)
	assert.NoError(t, err)
	defer repo.GracefulShutdown()

	tests := []struct {
		name          string
		fighter       model.Fighter
		expectedError error
	}{
		{
			name:          "Success",
			fighter:       *testFighter,
			expectedError: nil,
		},
		{
			name: "Error",
			fighter: model.Fighter{
				FighterId: 999888777,
			},
			expectedError: errors.New("no rows in result set"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fighterIO, err := repo.UpdateFighter(ctx, nil, tc.fighter)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, fighterIO, tc.fighter.FighterId)
			}
		})
	}
}

func TestUpdateFighterStats(t *testing.T) {
	initTestConfig()
	defer viper.Reset()

	ctx := context.Background()
	config := cfg.ViperTestPostgres()

	repo, err := New(ctx, config)
	assert.NoError(t, err)
	defer repo.GracefulShutdown()

	err = repo.UpdateFighterStats(ctx, nil, model.FighterStats{FighterId: testFighter.FighterId, TkdAccuracy: 50})
	assert.NoError(t, err)

	tx, err := repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	assert.NoError(t, err)

	err = repo.UpdateFighterStats(ctx, tx, model.FighterStats{FighterId: testFighter.FighterId, TkdAccuracy: 55})
	assert.NoError(t, err)

	err = tx.Commit(ctx)
	assert.NoError(t, err)
}

func initTestConfig() {
	env := os.Getenv("APP_ENV")

	if env == "local" {
		viper.SetConfigName("config")
		viper.AddConfigPath("../../../configs")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %s\n", err)
		}
	} else if env == "ci" {
		viper.Set("postgres.test.data_dir", os.Getenv("POSTGRES_DATA_DIR"))
		viper.Set("postgres.test.url", os.Getenv("POSTGRES_URL"))
		viper.Set("postgres.test.host", os.Getenv("POSTGRES_HOST"))
		viper.Set("postgres.test.port", os.Getenv("POSTGRES_PORT"))
		viper.Set("postgres.test.name", os.Getenv("POSTGRES_NAME"))
		viper.Set("postgres.test.user", os.Getenv("POSTGRES_USER"))
		viper.Set("postgres.test.password", os.Getenv("POSTGRES_PASSWORD"))
	}

	// viper.SetConfigName("config")
	// viper.AddConfigPath("../../../configs")
	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Fatalf("Error reading config file: %s\n", err)
	// }
}
