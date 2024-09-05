package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"pickfighter.com/fighters/internal/repository/psql"
	"pickfighter.com/fighters/pkg/cfg"
	fightersmodel "pickfighter.com/fighters/pkg/model"
	"pickfighter.com/pkg/model"
)

var testFighter = fightersmodel.Fighter{
	FighterId:      999888,
	Name:           "John Doe",
	NickName:       "The Phantom",
	Division:       1,
	Status:         "Active",
	Hometown:       "New York",
	TrainsAt:       "MMA Gym",
	FightingStyle:  "Boxing",
	Age:            30,
	Height:         72.5,
	Weight:         185.5,
	OctagonDebut:   "2022-03-01",
	DebutTimestamp: 1646102400,
	Reach:          74.0,
	LegReach:       40.0,
	Wins:           10,
	Loses:          2,
	Draw:           1,
	FighterUrl:     "https://example.com/john_doe",
	ImageUrl:       "https://example.com/john_doe.jpg",
	Stats: fightersmodel.FighterStats{
		StatId:               999888,
		FighterId:            999888,
		TotalSigStrLanded:    200,
		TotalSigStrAttempted: 300,
		StrAccuracy:          67,
		TotalTkdLanded:       50,
		TotalTkdAttempted:    100,
		TkdAccuracy:          50,
		SigStrLanded:         6.5,
		SigStrAbs:            4.0,
		SigStrDefense:        75,
		TakedownDefense:      80,
		TakedownAvg:          2.5,
		SubmissionAvg:        1.2,
		KnockdownAvg:         0.5,
		AvgFightTime:         "10:30",
		WinByKO:              5,
		WinBySub:             3,
		WinByDec:             2,
	},
}

func TestExecute(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"--version"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestBindViperPersistentFlag(t *testing.T) {
	cmd := &cobra.Command{}

	cmd.PersistentFlags().String("testPersistentFlag", "", "Test persistent flag")

	bindViperPersistentFlag(cmd, "testPersistentFlag", "testPersistentFlag")

	cmd.SetArgs([]string{"--testPersistentFlag=testPersistentViperVal"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "testPersistentViperVal", viper.GetString("testPersistentFlag"), "Expected Viper persistent flag to be bound")
}

func TestBindViperPersistentNilFlag(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)

	cmd := &cobra.Command{}

	bindViperPersistentFlag(cmd, "testPersistentFlag", "")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Failed to bind viper flag:", "Expected log message about failed to bind viper flag")

	log.SetOutput(os.Stderr)
}

func TestValidateServerArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected error
	}{
		{"ValidRoute", []string{model.FightersService}, nil},
		{"EmptyArgs", []string{}, errEmptyApiRoute},
		{"InvalidRoute", []string{"invalidRoute"}, fmt.Errorf("allowed routes are: %s", model.FightersService)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := validateServerArgs(cmd, tc.args)

			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestReadFighterData(t *testing.T) {
	fighters, err := ReadFighterData()

	assert.NoError(t, err)
	assert.True(t, len(fighters) > 0)
}

func TestWriteFighterData(t *testing.T) {
	initTestConfig()
	defer viper.Reset()

	ctx := context.Background()
	config := cfg.ViperTestPostgres()

	fighterData := []fightersmodel.Fighter{
		testFighter,
		getRandomFighter(),
		getRandomFighter(),
		getRandomFighter(),
	}

	err := WriteFighterData(ctx, fighterData, config)
	assert.NoError(t, err)
}

func TestDeleteFighterData(t *testing.T) {
	// initTestConfig()
	// defer viper.Reset()

	// ctx := context.Background()
	// config := cfg.ViperTestPostgres()

	// err := DeleteFighterData(ctx, config)
	// assert.NoError(t, err)
}

func TestCreateFighter(t *testing.T) {
	initTestConfig()
	defer viper.Reset()

	ctx := context.Background()
	config := cfg.ViperTestPostgres()

	repo, err := psql.New(ctx, config)
	defer repo.GracefulShutdown()
	assert.NoError(t, err)

	err = createFighter(ctx, repo, getRandomFighter())
	assert.NoError(t, err)

	err = createFighter(ctx, repo, testFighter)
	assert.Error(t, err)
}

func TestUpdateFighter(t *testing.T) {
	initTestConfig()
	defer viper.Reset()

	ctx := context.Background()
	config := cfg.ViperTestPostgres()

	repo, err := psql.New(ctx, config)
	defer repo.GracefulShutdown()
	assert.NoError(t, err)

	err = updateFighter(ctx, repo, testFighter)
	assert.Error(t, err)
}

func initTestConfig() {
	if os.Getenv("APP_ENV") != "ci" {
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	env := os.Getenv("APP_ENV")

	if env == "local" {
		viper.SetConfigName("config")
		viper.AddConfigPath("../configs")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %s\n", err)
		}
	} else if env == "ci" {
		viper.Set("postgres.test.data_dir", os.Getenv("POSTGRES_DATA_DIR"))
		viper.Set("postgres.test.url", os.Getenv("POSTGRES_URL"))
		viper.Set("postgres.test.host", "localhost")
		viper.Set("postgres.test.port", "5433")
		viper.Set("postgres.test.name", "fighters_db")
		viper.Set("postgres.test.user", "postgres")
		// viper.Set("postgres.test.host", os.Getenv("POSTGRES_HOST"))
		// viper.Set("postgres.test.port", os.Getenv("POSTGRES_PORT"))
		// viper.Set("postgres.test.name", os.Getenv("POSTGRES_NAME"))
		// viper.Set("postgres.test.user", os.Getenv("POSTGRES_USER"))
		// viper.Set("postgres.test.password", os.Getenv("POSTGRES_PASSWORD"))

		fmt.Println(viper.GetString("postgres.test.data_dir"), "TEST")
		fmt.Println(viper.GetString("postgres.test.url"), "TEST")
		fmt.Println(viper.GetString("postgres.test.host"), "TEST")
		fmt.Println(viper.GetString("postgres.test.port"), "TEST")
		fmt.Println(viper.GetString("postgres.test.name"), "TEST")
		fmt.Println(viper.GetString("postgres.test.user"), "TEST")
		fmt.Println(viper.GetString("postgres.test.password"), "TEST")
	}
}

func getRandomFighter() fightersmodel.Fighter {
	fighterId := int32(getRandomNum())
	statId := int32(getRandomNum())
	debut := getRandomNum() * 1000

	return fightersmodel.Fighter{
		FighterId:      fighterId,
		Name:           "Test Fighter " + randomString(5),
		NickName:       "The Phantom",
		Division:       1,
		Status:         "Active",
		Hometown:       "New York",
		TrainsAt:       "MMA Gym",
		FightingStyle:  "Boxing",
		Age:            30,
		Height:         72.5,
		Weight:         185.5,
		OctagonDebut:   "2022-03-01",
		DebutTimestamp: debut,
		Reach:          74.0,
		LegReach:       40.0,
		Wins:           10,
		Loses:          2,
		Draw:           1,
		FighterUrl:     randomString(32),
		ImageUrl:       randomString(32),
		Stats: fightersmodel.FighterStats{
			StatId:               statId,
			FighterId:            fighterId,
			TotalSigStrLanded:    200,
			TotalSigStrAttempted: 300,
			StrAccuracy:          67,
			TotalTkdLanded:       50,
			TotalTkdAttempted:    100,
			TkdAccuracy:          50,
			SigStrLanded:         6.5,
			SigStrAbs:            4.0,
			SigStrDefense:        75,
			TakedownDefense:      80,
			TakedownAvg:          2.5,
			SubmissionAvg:        1.2,
			KnockdownAvg:         0.5,
			AvgFightTime:         "10:30",
			WinByKO:              5,
			WinBySub:             3,
			WinByDec:             2,
		},
	}
}

func getRandomNum() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(99999)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
