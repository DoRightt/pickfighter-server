package repo

import (
	"context"
	"fightbettr.com/fb-server/pkg/logger"
	"fightbettr.com/fb-server/pkg/model"
	"fightbettr.com/fb-server/pkg/pgxs"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/assert.v1"
)

func TestTxCreateUser(t *testing.T) {
	updateTables()
	defer viper.Reset()

	ctx := context.Background()
	config := getConfig()

	db, err := pgxs.NewPool(context.Background(), logger.NewSugared(), config)
	require.NoError(t, err)

	authRepo := New(db)

	user := model.User{
		Name:  "test",
		Claim: "1",
	}

	tx, err := authRepo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	require.NoError(t, err)

	userId, err := authRepo.TxCreateUser(ctx, tx, user)
	require.NoError(t, err)

	assert.Equal(t, userId, int32(1))
}
