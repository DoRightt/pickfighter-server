package auth

import (
	"context"
	"time"

	internalErr "fightbettr.com/auth/pkg/errors"
	"fightbettr.com/auth/pkg/model"
	"fightbettr.com/auth/pkg/utils"
	"github.com/jackc/pgx/v5"
)

func (c *Controller) Register(ctx context.Context, req *model.RegisterRequest) (int32, error) {
	tx, err := c.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		c.Logger.Errorf("Unable to begin transaction: %s", err)
		cErr := internalErr.New(internalErr.Tx, err, 101)
		return 0, cErr
	}

	credentials, err := c.createUserCredentials(ctx, tx, req)
	if err != nil {
		c.Logger.Errorf("Error while user credentials creation: %s", err)
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		c.Logger.Errorf("Unable to commit transaction: %s", err)
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

func (c *Controller) RegisterConfirm(ctx context.Context, req *model.UserCredentialsRequest) (bool, error) {
	creds, err := c.repo.FindUserCredentials(ctx, model.UserCredentialsRequest{
		Token: req.Token,
	})
	if err != nil {
		// TODO handle error
		if err == pgx.ErrNoRows {
			return false, internalErr.New(internalErr.UserCredentialsToken, err, 401)
		} else {
			c.Logger.Errorf("Failed to get user credentials: %s", err)
			return false, internalErr.NewDefault(internalErr.UserCredentialsToken, 402)
		}
	}

	if time.Now().Unix() >= creds.TokenExpire {
		// TODO handle error
		// 	httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.TokenExpired,
		// 		fmt.Errorf("token expired, try to reset password"))
		return false, err
	}

	if err := c.repo.ConfirmCredentialsToken(ctx, nil, model.UserCredentialsRequest{
		UserId:    creds.UserId,
		Token:     creds.Token,
		TokenType: creds.TokenType,
	}); err != nil {
		// TODO handle error
		c.Logger.Errorf("Failed to update user credentials: %s", err)
		// 	httplib.ErrorResponseJSON(w, http.StatusInternalServerError, 111, err)
		return false, err
	}

	return true, nil
}

func (c *Controller) Login(ctx context.Context, req *model.AuthenticateRequest) (*model.AuthenticateResult, error) {
	creds, err := c.repo.FindUserCredentials(ctx, model.UserCredentialsRequest{
		Email: req.Email,
	})

	if err != nil {
		// TODO Handle errors
		// if err == pgx.ErrNoRows {
		// 	httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.UserCredentialsNotExists,
		// 		fmt.Errorf("%s", "User with specified login credentials not exists"))
		// 	return
		// } else {
		c.Logger.Errorf("Failed to get user credentials: %s", err)
		// httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.UserCredentials, err)
		// }
		return nil, err
	}

	if !creds.Active {
		// TODO Handle errors
		// httplib.ErrorResponseJSON(w, http.StatusForbidden, internalErr.UserCredentialsIsNotActive,
		// 	fmt.Errorf("%s", "User is not activated"))
		return nil, err
	}

	p := utils.GenerateSaltedHash(req.Password, creds.Salt)
	if p != creds.Password {
		// TODO Handle errors
		// httplib.ErrorResponseJSON(w, http.StatusBadRequest, 1, fmt.Errorf("%s", "Wrong password"))
		return nil, err
	}

	if req.RememberMe {
		req.ExpiresIn = 60 * 60 * 24 * 7
	} else {
		req.ExpiresIn = 60 * 60 * 24
	}

	token, err := c.createJWTToken(ctx, &creds, req)
	if err != nil {
		// TODO Handle error
		// c.Logger.Errorf("Unable to create session for google JWT: %s", err)
		return nil, err
	}

	return token, nil
}
