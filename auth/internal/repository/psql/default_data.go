package psql

import (
	"context"
	"math/rand"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"github.com/DoRightt/pickfighter-server/auth/pkg/model"
	"github.com/DoRightt/pickfighter-server/auth/pkg/utils"
	logs "github.com/DoRightt/pickfighter-server/pkg/logger"
	modeldefault "github.com/DoRightt/pickfighter-server/pkg/model"
)

type rootUserData struct {
	user        model.User
	credentials model.UserCredentials
}

func getRootUserData() rootUserData {
	user := model.User{
		Name:  "Hunter",
		Claim: "root_user",
		Flags: 1,
	}

	rand.Seed(time.Now().UnixNano())
	salt := utils.GetRandomString(modeldefault.DefaultSaltLength)
	passString := viper.GetString("root.password")
	password := utils.GenerateSaltedHash(passString, salt)
	email := viper.GetString("root.ownersEmail")

	token := utils.GenerateHashFromString(email + password + salt + user.Name)
	tokenExpire := time.Now().Unix() + 60*60*48

	userCredentials := model.UserCredentials{
		Email:       email,
		Password:    password,
		Salt:        salt,
		Token:       token,
		TokenType:   model.TokenConfirmation,
		TokenExpire: tokenExpire,
	}

	return rootUserData{
		user:        user,
		credentials: userCredentials,
	}
}

func (r *Repository) InitRootUser(ctx context.Context) error {
	userData := getRootUserData()
	existingAdminCredentials, err := r.FindUserCredentials(ctx, model.UserCredentialsRequest{
		Email: userData.credentials.Email,
	})

	if err == nil {
		user, err := r.FindUser(ctx, &model.UserRequest{UserId: existingAdminCredentials.UserId})
		if err != nil {
			return err
		}

		if user.Claim != userData.user.Claim {
			logs.Debugw("Current root profile does not have right claims",
				"current", user.Claim, "next", userData.user.Claim)
			user.Claim = userData.user.Claim
			if err := r.PatchUser(ctx, nil,
				user.UserId, "claim", userData.user.Claim,
			); err != nil {
				return err
			}
		}

		return nil
	}

	tx, err := r.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		logs.Errorw("Unable to begin transaction", "err", err)
		return err
	}

	userId, err := r.TxCreateUser(ctx, tx, userData.user)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		logs.Errorf("Failed to create user: %s", err)
		return err
	}

	userData.credentials.UserId = userId

	if err := r.TxNewAuthCredentials(ctx, tx, userData.credentials); err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		if err.(*pgconn.PgError).Code == "23505" {
			logs.Debug("duplicate error")
			return err
		} else {
			logs.Errorf("Failed to create user credentials during registration transaction: %s", err)
			return err
		}
	} else {
		logs.Info("Successfully created Root User credentials")
	}

	if err := r.ConfirmCredentialsToken(ctx, tx, model.UserCredentialsRequest{
		UserId:    userData.credentials.UserId,
		Token:     userData.credentials.Token,
		TokenType: userData.credentials.TokenType,
	}); err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		logs.Errorf("Unable to confirm user credentials: %s", err)
		return err
	} else {
		logs.Info("Successfully confirmed Captain credentials")
	}

	if err := tx.Commit(ctx); err != nil {
		logs.Errorf("Unable to commit transaction: %s", err)
		return err
	}

	return nil
}
