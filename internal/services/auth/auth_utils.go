package auth

import (
	"context"
	"net/http"
	"net/mail"
	internalErr "projects/fb-server/pkg/errors"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/utils"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

const (
	saltLength = 32
)

func (s service) createUserCredentials(ctx context.Context, tx pgx.Tx, req *model.RegisterRequest) (*model.UserCredentials, error) {
	if req.Email == "" {
		intErr := internalErr.New(internalErr.AuthFormEmailInvalid)
		return nil, httplib.NewApiErrFromInternalErr(intErr)
	}

	req.Email = strings.ToLower(req.Email)

	if _, err := mail.ParseAddress(req.Email); err != nil {
		intErr := internalErr.New(internalErr.AuthFormEmailInvalid)
		return nil, httplib.NewApiErrFromInternalErr(intErr)
	}

	if req.Password == "" || len(req.Password) < 6 {
		intErr := internalErr.New(internalErr.AuthFormPasswordInvalid)
		return nil, httplib.NewApiErrFromInternalErr(intErr)
	}

	user := model.User{
		Name:      req.Name,
		CreatedAt: time.Now().Unix(),
	}

	userId, err := s.Repo.TxCreateUser(ctx, tx, user)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			s.Logger.Errorf("Unable to rollback transaction: %s", txErr)
		}
		if err.(*pgconn.PgError).Code == pgerrcode.UniqueViolation {
			intErr := internalErr.New(internalErr.TxNotUnique)
			return nil, httplib.NewApiErrFromInternalErr(intErr)
		} else {
			intErr := internalErr.New(internalErr.TxUnknown)
			s.Logger.Errorf("Failed to create user during registration transaction: %s", err)
			return nil, httplib.NewApiErrFromInternalErr(intErr, http.StatusInternalServerError)
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

	if err := s.Repo.TxNewAuthCredentials(ctx, tx, userCredentials); err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			s.Logger.Errorf("Unable to rollback transaction: %s", txErr)
		}
		if err.(*pgconn.PgError).Code == pgerrcode.UniqueViolation {
			intErr := internalErr.New(internalErr.TxNotUnique)
			return nil, httplib.NewApiErrFromInternalErr(intErr)
		} else {
			intErr := internalErr.New(internalErr.TxUnknown)
			s.Logger.Errorf("Failed to create user during registration transaction: %s", err)
			return nil, httplib.NewApiErrFromInternalErr(intErr, http.StatusInternalServerError)
		}
	}

	return &userCredentials, nil
}
