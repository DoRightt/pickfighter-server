package repo

import (
	"context"
	"projects/fb-server/pkg/logger"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/pgxs"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestTxNewAuthCredentials(t *testing.T) {
	updateTables()
	defer viper.Reset()

	ctx := context.Background()
	config := getConfig()

	db, err := pgxs.NewPool(context.Background(), logger.NewSugared(), config)
	require.NoError(t, err)

	authRepo := New(db)

	userCredentials := model.UserCredentials{
		UserId:   int32(5),
		Email:    "test@gmail.com",
		Password: "test123qwerty",
		Salt:     "supersalt",
		Active:   false,
	}

	err = authRepo.TxNewAuthCredentials(ctx, nil, userCredentials)
	require.NoError(t, err)
}
