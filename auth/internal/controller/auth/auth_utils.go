package auth

import (
	"context"
	"net/mail"
	"strings"
	"time"

	internalErr "fightbettr.com/auth/pkg/errors"
	"fightbettr.com/auth/pkg/model"
	"fightbettr.com/auth/pkg/utils"
	logs "fightbettr.com/pkg/logger"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

const (
	saltLength = 32
)

// createUserCredentials creates user credentials during the user registration process.
// It validates the provided email and password, creates a new user, generates a salted hash
// of the password, and creates user credentials. If activation is required, it generates
// an email verification token. The created credentials are stored in the database.
func (c *Controller) createUserCredentials(ctx context.Context, tx pgx.Tx, req *model.RegisterRequest) (*model.UserCredentials, error) {
	if req.Email == "" {
		err := internalErr.NewDefault(internalErr.AuthFormEmailInvalid, 201)
		return nil, err
	}

	req.Email = strings.ToLower(req.Email)

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, internalErr.New(internalErr.AuthFormEmailInvalid, err, 202)
	}

	if req.Password == "" || len(req.Password) < 6 {
		return nil, internalErr.NewDefault(internalErr.AuthFormPasswordInvalid, 203)
	}

	user := model.User{
		Name:      req.Name,
		CreatedAt: time.Now().Unix(),
	}

	userId, err := c.repo.TxCreateUser(ctx, tx, user)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		pgErr, isPgError := err.(*pgconn.PgError)
		if isPgError && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, internalErr.New(internalErr.TxNotUnique, pgErr, 103)
		} else {
			logs.Errorf("Failed to create user during registration transaction: %s", err)
			return nil, internalErr.New(internalErr.TxUnknown, err, 104)
		}
	}

	salt := utils.GetRandomString(saltLength)
	password := utils.GenerateSaltedHash(req.Password, salt)

	activationDisabled := !viper.GetBool("auth.require_verification")

	userCredentials := model.UserCredentials{
		UserId:   userId,
		Email:    req.Email,
		Password: password,
		Salt:     salt,
		Active:   activationDisabled,
	}

	if !activationDisabled {
		userCredentials.Token = utils.GenerateHashFromString(req.Email + password + salt + req.Name)
		userCredentials.TokenExpire = time.Now().Unix() + 60*60*48
		userCredentials.TokenType = model.TokenConfirmation
	}

	if err := c.repo.TxNewAuthCredentials(ctx, tx, userCredentials); err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			logs.Errorf("Unable to rollback transaction: %s", txErr)
		}
		pgErr, isPgError := err.(*pgconn.PgError)
		if isPgError && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, internalErr.New(internalErr.TxNotUnique, pgErr, 105)
		} else {
			logs.Errorf("Failed to create user during registration transaction: %s", err)
			return nil, internalErr.New(internalErr.TxUnknown, err, 106)
		}
	}

	return &userCredentials, nil
}
