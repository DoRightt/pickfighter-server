package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/viper"
	authmodel "pickfighter.com/auth/pkg/model"
	eventmodel "pickfighter.com/events/pkg/model"
	fightersmodel "pickfighter.com/fighters/pkg/model"
	"pickfighter.com/pkg/httplib"
	"pickfighter.com/pkg/model"
	"pickfighter.com/pkg/utils"

	internalErr "pickfighter.com/pickfighter/pkg/errors"
)

// * * * * * Fighters Handlers * * * * *

// GetFighters handles HTTP requests to retrieve fighters based on status.
func (h *Handler) GetFighters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := utils.Capitalize(r.FormValue("status"))

	fighters, err := h.ctrl.SearchFighters(ctx, fightersmodel.FightersRequest{Status: status})
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	httplib.ResponseJSON(w, httplib.ListResult{
		Results: fighters,
		Count:   int32(len(fighters)),
	})
}

// * * * * * Auth Handlers * * * * *

// Register handles the registration of a new user.
// It expects a JSON request with user details, including name, email, password, and terms agreement.
// Upon successful registration, it initiates a confirmation email and returns the user's ID.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req authmodel.RegisterRequest
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.AuthDecode, err)
	}

	if !req.TermsOk {
		httplib.ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			internalErr.AuthForm,
			fmt.Errorf("you must accept terms and contiditons 'terms_ok' set to true"),
		)
		return
	}

	credentials, err := h.ctrl.Register(ctx, &req)
	if err != nil {
		credErr := err.(httplib.ApiError)
		httplib.ErrorResponseJSON(w, credErr.HttpStatus, credErr.ErrorCode, err)
		return
	}

	result := httplib.SuccessfulResult()
	result.Id = credentials.UserId

	httplib.ResponseJSON(w, result)
}

// ConfirmRegistration handles the confirmation of user registration by validating the provided token.
// Users receive a confirmation token upon successful registration, and this endpoint is used to confirm
// and activate their accounts
func (h *Handler) ConfirmRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.FormValue("token")
	if token == "" {
		httplib.ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			internalErr.QueryParamsToken,
			fmt.Errorf("query parameter 'token' should be specified"),
		)
		return
	}

	_, err := h.ctrl.ConfirmRegistration(ctx, token)
	if err != nil {
		// TODO handle errors from service
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, 111, err)
		return
	}

	httplib.ResponseJSON(w, httplib.SuccessfulResult())
}

// Login handles the user login process, authenticating the user based on the provided credentials.
// It validates the email or username and password, checks user activation status,
// generates a JWT token for the authenticated user, and sets an authentication cookie.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	var req authmodel.AuthenticateRequest
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.AuthDecode, err)
		return
	}

	if req.Email == "" {
		httplib.ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			internalErr.AuthFormEmailEmpty,
			fmt.Errorf("%s", "Empty 'email' or 'username'"),
		)
		return
	}

	if req.Password == "" {
		httplib.ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			internalErr.AuthFormPasswordInvalid,
			fmt.Errorf("%s", "Empty 'password'"),
		)
		return
	}

	req.Email = strings.ToLower(req.Email)
	req.UserAgent = r.UserAgent()
	// TODO
	// req.IpAddress = r.Header.Get(ipaddr.CFConnectingIp)

	token, err := h.ctrl.Login(ctx, &req)
	if err != nil {
		// TODO handle error
		httplib.ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			internalErr.Auth,
			err,
		)
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

// Logout handles the user logout process by setting an expired cookie.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
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

// ResetPassword handles the process of resetting a user's password.
// It expects a JSON request containing the user's email address.
// If the email is valid and associated with an existing user, a reset token is generated,
// and an email containing the reset link is sent to the user.
// The reset token is also stored in the database for verification during the password reset process.
// A successful response is returned if the email exists, and the reset process is initiated.
// In case of errors, appropriate error responses are sent with details in the response body.
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// TODO add other fields except email
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req authmodel.ResetPasswordRequest
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.AuthDecode, err)
		return
	}

	noEmail := len(req.Email) < 1 || req.Email == " "
	if noEmail {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.AuthFormEmailEmpty,
			fmt.Errorf("%s", "Empty 'email'"))
		return
	}

	req.Email = strings.ToLower(req.Email)

	_, err := h.ctrl.ResetPassword(ctx, &req)
	if err != nil {
		// TODO handle errors from service
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, 111, err)
		return
	}

	httplib.ResponseJSON(w, httplib.SuccessfulResult())
}

// RecoverPassword handles the process of recovering a user's password based on a provided reset token.
// It expects a JSON request containing the reset token, new password, and confirmation password.
// If the token is valid, the password is updated, and the token is marked as used.
// The response includes a successful result if the password recovery process is completed.
// In case of errors, appropriate error responses are sent with details in the response body.
func (h *Handler) RecoverPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req authmodel.RecoverPasswordRequest
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.AuthDecode, err)
		return
	}

	if len(req.Token) < 2 || req.Token == " " {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Token, fmt.Errorf("empty 'token'"))
		return
	}

	noPassword := len(req.Password) < 6
	noConfirm := len(req.ConfirmPassword) < 6

	if noPassword && noConfirm {
		httplib.ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			internalErr.AuthFormPasswordInvalid,
			fmt.Errorf("empty body 'password'"),
		)
		return
	}

	if noConfirm {
		httplib.ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			internalErr.Auth,
			fmt.Errorf("empty body 'confirm_password'"),
		)
		return
	}

	if req.Password != req.ConfirmPassword {
		httplib.ErrorResponseJSON(
			w,
			http.StatusBadRequest,
			internalErr.AuthFormPasswordsMismatch,
			fmt.Errorf("password are not equal"),
		)
		return
	}

	_, err := h.ctrl.PasswordRecover(ctx, &req)
	if err != nil {
		// TODO handle errors from service
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, 111, err)
		return
	}

	httplib.ResponseJSON(w, httplib.SuccessfulResult())
}

// GetCurrentUser retrieves information about the currently authenticated user.
// It extracts the user ID from the request context, queries the database for the user's details,
// and responds with a JSON representation of the user.
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := h.ctrl.GetCurrentUser(ctx)
	if err != nil {
		// TODO handle errors from service
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, 111, err)
		return
	}

	httplib.ResponseJSON(w, authmodel.UserResult{
		User: *user,
	})
}

// * * * * * Event Handlers * * * * *

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req eventmodel.EventRequest
	if err := decoder.Decode(&req); err != nil {
		// TODO handle errors from service
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Events, err)
	}

	event, err := h.ctrl.CreateEvent(ctx, &req)
	if err != nil {
		// TODO handle errors from service
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Events, err)
	}

	result := httplib.SuccessfulResult()
	result.Id = event.EventId

	httplib.ResponseJSON(w, result)
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.ctrl.SearchEvents(ctx)
	if err != nil {
		// TODO handle errors from service
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Events, err)
	}

	httplib.ResponseJSON(w, httplib.ListResult{
		Results: res.Events,
		Count:   res.Count,
	})
}

func (h *Handler) CreateBet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req eventmodel.Bet
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Events, err)
	}

	token, ok := ctx.Value(model.ContextJWTPointer).(jwt.Token)
	if !ok {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, 322,
			fmt.Errorf("unable to find request context token"))
		return
	}

	userId, ok := token.Get(string(model.ContextUserId))
	if !ok {
		httplib.ErrorResponseJSON(w, http.StatusUnauthorized, http.StatusUnauthorized,
			fmt.Errorf("illegal token, user id must be specified"))
		return
	}

	req.UserId = int32(userId.(float64))

	betId, err := h.ctrl.CreateBet(ctx, &req)
	if err != nil {
		// TODO handle errors from service
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Bets, err)
	}

	result := httplib.SuccessfulResult()
	result.Id = betId

	httplib.ResponseJSON(w, result)
}

func (h *Handler) GetBets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, ok := ctx.Value(model.ContextJWTPointer).(jwt.Token)
	if !ok {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, 323,
			fmt.Errorf("unable to find request context token"))
		return
	}

	userId, ok := token.Get(string(model.ContextUserId))
	if !ok {
		httplib.ErrorResponseJSON(w, http.StatusUnauthorized, http.StatusUnauthorized,
			fmt.Errorf("illegal token, user id must be specified"))
		return
	}

	id, ok := userId.(float64)
	if !ok {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, 324,
			fmt.Errorf("unable to convert user id"))
		return
	}

	resp, err := h.ctrl.SearchBets(ctx, int32(id))
	if err != nil {
		// TODO handle errors from service
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Bets, err)
		return
	}

	httplib.ResponseJSON(w, httplib.ListResult{
		Results: resp.Bets,
		Count:   resp.Count,
	})
}

func (h *Handler) AddResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req eventmodel.FightResultRequest
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Events, err)
	}

	id, err := h.ctrl.SetResult(ctx, &req)
	if err != nil {
		// TODO handle errors from service
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Bets, err)
	}

	result := httplib.SuccessfulResult()
	result.Id = id

	httplib.ResponseJSON(w, result)
}
