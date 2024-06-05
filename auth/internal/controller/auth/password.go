package auth

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"fightbettr.com/auth/pkg/model"
	"fightbettr.com/auth/pkg/utils"
	"github.com/jackc/pgx/v5"
)

func (c *Controller) PasswordReset(ctx context.Context, req *model.ResetPasswordRequest) (bool, error) {
	credentials, err := c.repo.FindUserCredentials(ctx, model.UserCredentialsRequest{
		Email: req.Email,
	})
	if err != nil {
		// TODO handle errors
		// if err == pgx.ErrNoRows {
		// 	httplib.ErrorResponseJSON(w, http.StatusNotFound, http.StatusNotFound, err)
		// 	return
		// } else {
		c.Logger.Errorf("Failed to find user credentials: %s", err)
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.UserCredentials, err)
		// return
		// }
		return false, err
	}

	user, err := c.repo.FindUser(ctx, &model.UserRequest{
		UserId: credentials.UserId,
	})
	if err != nil {
		// TODO handle error
		c.Logger.Errorf("Failed to find user: %s", err)
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.Profile, err)
		return false, err
	}

	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		// TODO handle error
		c.Logger.Errorf("Failed to create registration transaction: %s", err)
		// httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Tx, err)
		return false, err
	}

	rn := rand.New(rand.NewSource(time.Now().UnixNano()))
	salt := rn.Int()

	token := utils.GenerateHashFromString(fmt.Sprintf("%s:%s:%d", req.Email, time.Now(), +salt))
	tokenExpire := time.Now().Unix() + 60*60*48
	credentials.TokenType = model.TokenResetPassword
	credentials.Token = token
	credentials.TokenExpire = tokenExpire

	if err := c.repo.ResetPassword(ctx, &credentials); err != nil {
		// TODO handle error
		c.Logger.Errorf("Failed to reset user credentials: %s", err)
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.TxCommit, err)
		return false, err
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		// TODO handle error
		c.Logger.Errorf("Failed to commit registration transaction: %s", txErr)
		// httplib.ErrorResponseJSON(w, http.StatusBadRequest, 11, txErr)
		return false, err
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
		// TODO handle service errors
		// if err == pgx.ErrNoRows {
		// 	httplib.ErrorResponseJSON(w, http.StatusNotFound, http.StatusNotFound, err)
		// 	return
		// } else {
		c.Logger.Errorf("Failed to find user credentials: %s", err)
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.UserCredentials, err)
		// return
		// }
		return false, err
	}

	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		// TODO handle service errors
		c.Logger.Errorf("Failed to create registration transaction: %s", err)
		// httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Tx, err)
		return false, err
	}

	salt := utils.GetRandomString(saltLength)
	password := utils.GenerateSaltedHash(req.Password, salt)

	credentials.Password = password
	credentials.Salt = salt

	if err := c.repo.ConfirmCredentialsToken(ctx, tx, model.UserCredentialsRequest{
		UserId: credentials.UserId,
	}); err != nil {
		// TODO handle service errors
		c.Logger.Errorf("Failed to reset user credentials: %s", err)
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.UserCredentials, err)
		return false, err
	}

	if err := c.repo.UpdatePassword(ctx, tx, credentials); err != nil {
		// TODO handle service errors
		c.Logger.Errorf("Failed to update user password: %s", err)
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.UserCredentialsReset, err)
		return false, err
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		// TODO handle service errors
		c.Logger.Errorf("Failed to commit registration transaction: %s", txErr)
		// httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.TxCommit, txErr)
		return false, err
	}

	return true, nil
}
