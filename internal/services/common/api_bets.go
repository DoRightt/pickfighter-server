package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	internalErr "projects/fb-server/pkg/errors"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

// GetBets retrieves and returns the bets associated with the authenticated user. It extracts the
// user ID from the JWT token in the request context, calls the SearchBetsCount and SearchBets methods
// to get the count and list of bets, and responds with the results. If any error occurs during the
// process, it responds with an appropriate API error along with the HTTP status code. If there are no
// bets, it responds with an empty list result.
func (s *service) GetBets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, ok := ctx.Value(model.ContextJWTPointer).(jwt.Token)
	if !ok {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, 320,
			fmt.Errorf("unable to find request context token"))
		return
	}

	userId, ok := token.Get(string(model.ContextUserId))
	if !ok {
		httplib.ErrorResponseJSON(w, http.StatusUnauthorized, http.StatusUnauthorized,
			fmt.Errorf("illegal token, user id must be specified"))
		return
	}

	count, err := s.Repo.SearchBetsCount(ctx, int32(userId.(float64)))
	if err != nil {
		s.Logger.Errorf("Failed to get events count: %s", err)
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.CountBets, err)
		return
	}
	if count == 0 {
		httplib.ResponseJSON(w, httplib.ListResult{})
		return
	}

	bets, err := s.Repo.SearchBets(ctx, int32(userId.(float64)))
	if err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Bets, err)
	}

	httplib.ResponseJSON(w, httplib.ListResult{
		Results: bets,
		Count:   count,
	})

}

// CreateBet handles the creation of a new bet. It decodes the JSON request body into the
// model.Bet struct, retrieves the user ID from the JWT token in the request context, and calls
// the CreateBet method to create the bet. If any error occurs during the process, it responds with
// an appropriate API error along with the HTTP status code. If the creation is successful, it
// responds with a successful result and the ID of the created bet.
func (s *service) CreateBet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req model.Bet
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Events, err)
	}

	token, ok := ctx.Value(model.ContextJWTPointer).(jwt.Token)
	if !ok {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, 320,
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

	betId, err := s.Repo.CreateBet(ctx, &req)
	if err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Bets, err)
	}

	result := httplib.SuccessfulResult()
	result.Id = betId

	httplib.ResponseJSON(w, result)
}
