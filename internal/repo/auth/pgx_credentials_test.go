package repo

import (
	"context"
	"projects/fb-server/pkg/logger"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/pgxs"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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

func TestFindUserCredentials(t *testing.T) {
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

	tests := []struct {
		name               string
		credentialsRequest model.UserCredentialsRequest
	}{
		{
			name: "Request with Email",
			credentialsRequest: model.UserCredentialsRequest{
				UserId: userCredentials.UserId,
				Email:  userCredentials.Email,
			},
		},
		{
			name: "Request with UserId",
			credentialsRequest: model.UserCredentialsRequest{
				UserId: userCredentials.UserId,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			creds, err := authRepo.FindUserCredentials(ctx, tc.credentialsRequest)
			require.NoError(t, err)

			assert.Equal(t, creds.Email, userCredentials.Email)
			assert.Equal(t, creds.UserId, userCredentials.UserId)
		})
	}
}

func TestConfirmCredentialsToken(t *testing.T) {
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

	tests := []struct {
		name               string
		credentialsRequest model.UserCredentialsRequest
	}{
		{
			name: "Success",
			credentialsRequest: model.UserCredentialsRequest{
				UserId: userCredentials.UserId,
				Email:  userCredentials.Email,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := authRepo.ConfirmCredentialsToken(ctx, nil, tc.credentialsRequest)
			require.NoError(t, err)

			creds, err := authRepo.FindUserCredentials(ctx, tc.credentialsRequest)
			require.NoError(t, err)

			assert.True(t, creds.Active)
		})
	}
}

func TestResetPassword(t *testing.T) {
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

	tests := []struct {
		name               string
		credentials        model.UserCredentials
		credentialsRequest model.UserCredentialsRequest
	}{
		{
			name: "Success",
			credentials: model.UserCredentials{
				UserId: userCredentials.UserId,
				Email:  userCredentials.Email,
			},
			credentialsRequest: model.UserCredentialsRequest{
				UserId: userCredentials.UserId,
				Email:  userCredentials.Email,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := authRepo.ConfirmCredentialsToken(ctx, nil, tc.credentialsRequest)
			require.NoError(t, err)

			err = authRepo.ResetPassword(ctx, &tc.credentials)
			require.NoError(t, err)

			creds, err := authRepo.FindUserCredentials(ctx, tc.credentialsRequest)
			require.NoError(t, err)

			assert.False(t, creds.Active)
		})
	}
}

func TestUpdatePassword(t *testing.T) {
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

	tests := []struct {
		name               string
		credentials        model.UserCredentials
		credentialsRequest model.UserCredentialsRequest
	}{
		{
			name: "Success",
			credentials: model.UserCredentials{
				UserId:   userCredentials.UserId,
				Email:    userCredentials.Email,
				Password: "somenewpassword",
			},
			credentialsRequest: model.UserCredentialsRequest{
				UserId: userCredentials.UserId,
				Email:  userCredentials.Email,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := authRepo.ConfirmCredentialsToken(ctx, nil, tc.credentialsRequest)
			require.NoError(t, err)

			err = authRepo.UpdatePassword(ctx, nil, tc.credentials)
			require.NoError(t, err)

			creds, err := authRepo.FindUserCredentials(ctx, tc.credentialsRequest)
			require.NoError(t, err)

			assert.Equal(t, creds.Password, tc.credentials.Password)
		})
	}
}
