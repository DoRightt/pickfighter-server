package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	internalErr "projects/fb-server/pkg/errors"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/utils"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

// Register is a handler method for /register path
func (s *service) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req model.RegisterRequest
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.AuthDecode, err)
	}

	if !req.TermsOk {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.AuthForm,
			fmt.Errorf("you must accept terms and contiditons 'terms_ok' set to true"))
		return
	}

	tx, err := s.Repo.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		s.Logger.Errorf("Unable to begin transaction: %s", err)
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Tx, err)
	}

	credentials, err := s.createUserCredentials(ctx, tx, &req)
	if err != nil {
		credErr := err.(httplib.ApiError)
		httplib.ErrorResponseJSON(w, credErr.HttpStatus, credErr.ErrorCode, err)
		return
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		s.Logger.Errorf("Unable to commit transaction: %s", txErr)
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.TxCommit, txErr)
		return
	}

	go s.HandleEmailEvent(ctx, &model.EmailData{
		Subject: model.EmailRegistration,
		Recipient: model.EmailAddrSpec{
			Email: req.Email,
			Name:  req.Name,
		},
		Token: credentials.Token,
	})

	result := httplib.SuccessfulResult()
	result.Id = credentials.UserId

	httplib.ResponseJSON(w, result)
}

func (s *service) ConfirmRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.FormValue("token")
	if token == "" {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.QueryParamsToken,
			fmt.Errorf("query parameter 'token' should be specified"))
		return
	}

	creds, err := s.Repo.FindUserCredentials(ctx, model.UserCredentialsRequest{
		Token: token,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.UserCredentialsToken, err)
		} else {
			s.Logger.Errorf("Failed to get user credentials: %s", err)
			httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.UserCredentials, err)
		}
		return
	}

	if time.Now().Unix() >= creds.TokenExpire {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.TokenExpired,
			fmt.Errorf("token expired, try to reset password"))
		return
	}

	if err := s.Repo.ConfirmCredentialsToken(ctx, nil, model.UserCredentialsRequest{
		UserId:    creds.UserId,
		Token:     creds.Token,
		TokenType: creds.TokenType,
	}); err != nil {
		s.Logger.Errorf("Failed to update user credentials: %s", err)
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, 111, err)
		return
	}

	httplib.ResponseJSON(w, httplib.SuccessfulResult())
}

func (s *service) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	var req model.AuthenticateRequest
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.AuthDecode, err)
		return
	}

	if req.Email == "" {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.AuthFormEmailEmpty, fmt.Errorf("%s", "Empty 'email' or 'username'"))
		return
	}

	if req.Password == "" {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.AuthFormPasswordInvalid, fmt.Errorf("%s", "Empty 'password'"))
		return
	}

	req.Email = strings.ToLower(req.Email)

	creds, err := s.Repo.FindUserCredentials(ctx, model.UserCredentialsRequest{
		Email: req.Email,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.UserCredentialsNotExists,
				fmt.Errorf("%s", "User with specified login credentials not exists"))
			return
		} else {
			s.Logger.Errorf("Failed to get user credentials: %s", err)
			httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.UserCredentials, err)
			return
		}
	}
	if !creds.Active {
		httplib.ErrorResponseJSON(w, http.StatusForbidden, internalErr.UserCredentialsIsNotActive,
			fmt.Errorf("%s", "User is not activated"))
		return
	}

	p := utils.GenerateSaltedHash(req.Password, creds.Salt)
	if p != creds.Password {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, 1, fmt.Errorf("%s", "Wrong password"))
		return
	}

	req.UserAgent = r.UserAgent()
	// TODO
	// req.IpAddress = r.Header.Get(ipaddr.CFConnectingIp)

	if req.RememberMe {
		req.ExpiresIn = 60 * 60 * 24 * 7
	} else {
		req.ExpiresIn = 60 * 60 * 24
	}

	token, err := s.createJWTToken(ctx, &creds, req)
	if err != nil {
		s.Logger.Errorf("Unable to create session for google JWT: %s", err)
		return
	}

	authCookieName := viper.GetString("auth.cookie_name")
	http.SetCookie(w, &http.Cookie{
		Name:    authCookieName,
		Value:   token.AccessToken,
		Expires: token.ExpirationTime,
		Path:    "/",
	})

	result := httplib.SuccessfulResultMap()
	result["token_id"] = token.TokenId
	result["access_token"] = token.AccessToken
	result["expires_at"] = token.ExpirationTime
	httplib.ResponseJSON(w, result)
}

func (s *service) Logout(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// token, ok := ctx.Value(model.ContextJWTPointer).(jwt.Token)
	// if !ok {
	// 	httplib.ErrorResponseJSON(w, http.StatusBadRequest, 320,
	// 		fmt.Errorf("unable to find request context token"))
	// 	return
	// }

	// * * * * *

	http.SetCookie(w, &http.Cookie{
		Name:    viper.GetString("auth.cookie_name"),
		Value:   "",
		Expires: time.Now().Add(1 * time.Second),
		Path:    "/",
	})

	httplib.ResponseJSON(w, httplib.SuccessfulResultMap())
}
