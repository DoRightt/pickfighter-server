package repo

import (
	"context"
	"fightbettr.com/fb-server/pkg/logger"
	"fightbettr.com/fb-server/pkg/model"
	"fightbettr.com/fb-server/pkg/pgxs"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindUser(t *testing.T) {
	// TODO
}

func TestSearchUsers(t *testing.T) {
	// TODO
}

func TestPerformUsersRequestQuery(t *testing.T) {
	tests := []struct {
		name           string
		userReq        *model.UsersRequest
		expectedResult []string
	}{
		{
			name:           "No body",
			userReq:        nil,
			expectedResult: nil,
		},
		{
			name: "Single user id",
			userReq: &model.UsersRequest{
				UserIds: []int32{1},
			},
			expectedResult: []string{"u.user_id = 1"},
		},
		{
			name: "Multiple user ids",
			userReq: &model.UsersRequest{
				UserIds: []int32{1, 2, 3},
			},
			expectedResult: []string{"u.user_id IN (1,2,3)"},
		},
		{
			name: "With name",
			userReq: &model.UsersRequest{
				Name:    "test",
				UserIds: []int32{1, 2, 3},
			},
			expectedResult: []string{"u.user_id IN (1,2,3)", "u.name ILIKE '%test%'"},
		},
		{
			name: "With email",
			userReq: &model.UsersRequest{
				Name:    "test",
				Email:   "test@mail.com",
				UserIds: []int32{1, 2, 3},
			},
			expectedResult: []string{"u.user_id IN (1,2,3)", "u.name ILIKE '%test%'", "u.public_email ILIKE '%test@mail.com%'"},
		},
		{
			name: "With CreatedFrom",
			userReq: &model.UsersRequest{
				Name:    "test",
				Email:   "test@mail.com",
				UserIds: []int32{1, 2, 3},
				ListRequest: model.ListRequest{
					CreatedFrom: 123,
				},
			},
			expectedResult: []string{
				"u.user_id IN (1,2,3)",
				"u.name ILIKE '%test%'",
				"u.public_email ILIKE '%test@mail.com%'",
				"u.created_at > '123'",
			},
		},
		{
			name: "With CreatedUntil",
			userReq: &model.UsersRequest{
				Name:    "test",
				Email:   "test@mail.com",
				UserIds: []int32{1, 2, 3},
				ListRequest: model.ListRequest{
					CreatedFrom:  123,
					CreatedUntil: 123,
				},
			},
			expectedResult: []string{
				"u.user_id IN (1,2,3)",
				"u.name ILIKE '%test%'",
				"u.public_email ILIKE '%test@mail.com%'",
				"u.created_at > '123'",
				"u.created_at < '123'",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer viper.Reset()
			initTestConfig()

			config := &pgxs.Config{
				DataDir:  viper.GetString("postgres.test.data_dir"),
				DbUri:    viper.GetString("postgres.test.url"),
				Host:     viper.GetString("postgres.test.host"),
				Port:     viper.GetString("postgres.test.port"),
				Name:     viper.GetString("postgres.test.name"),
				User:     viper.GetString("postgres.test.user"),
				Password: viper.GetString("postgres.test.password"),
			}
			db, err := pgxs.NewPool(context.Background(), logger.NewSugared(), config)
			require.NoError(t, err)

			authRepo := New(db)

			result := authRepo.PerformUsersRequestQuery(tc.userReq)

			isEqual := reflect.DeepEqual(result, tc.expectedResult)

			assert.True(t, isEqual, "Expected %v, got %v", tc.expectedResult, result)
		})
	}
}
