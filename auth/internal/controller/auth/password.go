package auth

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	internalErr "fightbettr.com/auth/pkg/errors"
	"fightbettr.com/auth/pkg/model"
	"fightbettr.com/auth/pkg/utils"
	"github.com/jackc/pgx/v5"
)

func (c *Controller) PasswordReset(ctx context.Context, req *model.ResetPasswordRequest) (bool, error) {
	credentials, err := c.repo.FindUserCredentials(ctx, model.UserCredentialsRequest{
		Email: req.Email,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			// not found error
			return false, internalErr.New(internalErr.UserCredentials, err, 407)
		} else {
			// internal error
			c.Logger.Errorf("Failed to find user credentials: %s", err)
			return false, internalErr.New(internalErr.UserCredentials, err, 408)
		}
	}

	user, err := c.repo.FindUser(ctx, &model.UserRequest{
		UserId: credentials.UserId,
	})
	if err != nil {
		// internal error
		c.Logger.Errorf("Failed to find user: %s", err)
		return false, internalErr.New(internalErr.Profile, err, 501)
	}

	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		// bad request error
		c.Logger.Errorf("Failed to create registration transaction: %s", err)
		return false, internalErr.New(internalErr.Tx, err, 107)
	}

	rn := rand.New(rand.NewSource(time.Now().UnixNano()))
	salt := rn.Int()

	token := utils.GenerateHashFromString(fmt.Sprintf("%s:%s:%d", req.Email, time.Now(), +salt))
	tokenExpire := time.Now().Unix() + 60*60*48
	credentials.TokenType = model.TokenResetPassword
	credentials.Token = token
	credentials.TokenExpire = tokenExpire

	if err := c.repo.ResetPassword(ctx, &credentials); err != nil {
		// internal error
		c.Logger.Errorf("Failed to reset user credentials: %s", err)
		return false, internalErr.New(internalErr.TxCommit, err, 108)
	}

	if err := tx.Commit(ctx); err != nil {
		// bad request error
		c.Logger.Errorf("Failed to commit registration transaction: %s", err)
		return false, internalErr.New(internalErr.TxCommit, err, 109)
	}

	// TODO
	go c.HandleEmailEvent(ctx, &model.EmailData{
		Subject: model.EmailResetPassword,
		Recipient: model.EmailAddrSpec{
			Email: credentials.Email,
			Name:  user.Name,
		},
		Token: credentials.Token,
	})

	return true, nil
}

func (c *Controller) PasswordRecover(ctx context.Context, req *model.RecoverPasswordRequest) (bool, error) {
	credentials, err := c.repo.FindUserCredentials(ctx, model.UserCredentialsRequest{
		Token: req.Token,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			// not found error
			return false, internalErr.New(internalErr.UserCredentialsToken, err, 409)
		} else {
			// internal error
			c.Logger.Errorf("Failed to find user credentials: %s", err)
			return false, internalErr.New(internalErr.UserCredentials, err, 410)
		}
	}

	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		// bad request error
		c.Logger.Errorf("Failed to create registration transaction: %s", err)
		return false, internalErr.New(internalErr.Tx, err, 110)
	}

	salt := utils.GetRandomString(saltLength)
	password := utils.GenerateSaltedHash(req.Password, salt)

	credentials.Password = password
	credentials.Salt = salt

	if err := c.repo.ConfirmCredentialsToken(ctx, tx, model.UserCredentialsRequest{
		UserId: credentials.UserId,
	}); err != nil {
		// internal error
		c.Logger.Errorf("Failed to reset user credentials: %s", err)
		return false, internalErr.New(internalErr.UserCredentials, err, 411)
	}

	if err := c.repo.UpdatePassword(ctx, tx, credentials); err != nil {
		// internal error
		c.Logger.Errorf("Failed to update user password: %s", err)
		return false, internalErr.New(internalErr.UserCredentialsReset, err, 412)
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		// bad request error
		c.Logger.Errorf("Failed to commit registration transaction: %s", txErr)
		return false, internalErr.New(internalErr.TxCommit, err, 111)
	}

	return true, nil
}
