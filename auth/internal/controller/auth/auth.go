package auth

import (
	"context"
	"time"

	internalErr "fightbettr.com/auth/pkg/errors"
	"fightbettr.com/auth/pkg/model"
	"fightbettr.com/auth/pkg/utils"
	logs "fightbettr.com/pkg/logger"
	"github.com/jackc/pgx/v5"
)

// Register creates user credentials in a transactional context,
// commits the transaction upon success, and asynchronously sends an email confirmation.
// It returns the user ID upon successful registration.
func (c *Controller) Register(ctx context.Context, req *model.RegisterRequest) (int32, error) {
	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		logs.Errorf("Unable to begin transaction: %s", err)
		cErr := internalErr.New(internalErr.Tx, err, 101)
		return 0, cErr
	}

	credentials, err := c.createUserCredentials(ctx, tx, req)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		logs.Errorf("Error while user credentials creation: %s", err)
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		logs.Errorf("Unable to commit transaction: %s", err)
		cErr := internalErr.New(internalErr.TxCommit, err, 102)
		return 0, cErr
	}

	// TODO
	go c.HandleEmailEvent(ctx, &model.EmailData{
		Subject: model.EmailRegistration,
		Recipient: model.EmailAddrSpec{
			Email: req.Email,
			Name:  req.Name,
		},
		Token: credentials.Token,
	})

	return credentials.UserId, nil
}

// RegisterConfirm verifies the user credentials token,
// checks if the token is expired, and confirms the user's registration.
// It returns true if the token is valid and registration is confirmed successfully.
func (c *Controller) RegisterConfirm(ctx context.Context, req *model.UserCredentialsRequest) (bool, error) {
	creds, err := c.repo.FindUserCredentials(ctx, model.UserCredentialsRequest{
		Token: req.Token,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, internalErr.New(internalErr.UserCredentialsToken, err, 401)
		} else {
			logs.Errorf("Failed to get user credentials: %s", err)
			return false, internalErr.NewDefault(internalErr.UserCredentials, 402)
		}
	}

	if time.Now().Unix() >= creds.TokenExpire {
		return false, internalErr.NewDefault(internalErr.TokenExpired, 601)
	}

	if err := c.repo.ConfirmCredentialsToken(ctx, nil, model.UserCredentialsRequest{
		UserId:    creds.UserId,
		Token:     creds.Token,
		TokenType: creds.TokenType,
	}); err != nil {
		logs.Errorf("Failed to update user credentials: %s", err)
		return false, internalErr.New(internalErr.UserCredentialsUpdate, err, 403)
	}

	return true, nil
}

// Login verifies user credentials by email and password,
// generates a JWT token for authentication, and returns it.
// Returns an error if credentials are invalid or token generation fails.
func (c *Controller) Login(ctx context.Context, req *model.AuthenticateRequest) (*model.AuthenticateResult, error) {
	creds, err := c.repo.FindUserCredentials(ctx, model.UserCredentialsRequest{
		Email: req.Email,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, internalErr.New(internalErr.UserCredentialsNotExists, err, 404)
		} else {
			logs.Errorf("Failed to get user credentials: %s", err)
			return nil, internalErr.New(internalErr.UserCredentials, err, 405)
		}

	}

	if !creds.Active {
		return nil, internalErr.NewDefault(internalErr.UserCredentialsIsNotActive, 406)
	}

	p := utils.GenerateSaltedHash(req.Password, creds.Salt)
	if p != creds.Password {
		return nil, internalErr.NewDefault(internalErr.AuthFormPasswordWrong, 204)
	}

	if req.RememberMe {
		req.ExpiresIn = 60 * 60 * 24 * 7
	} else {
		req.ExpiresIn = 60 * 60 * 24
	}

	token, err := c.createJWTToken(ctx, &creds, req)
	if err != nil {
		logs.Errorf("Unable to create session for google JWT: %s", err)
		return nil, internalErr.New(internalErr.Token, err, 602)
	}

	return token, nil
}
