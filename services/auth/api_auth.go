package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	internalErr "projects/fb-server/errors"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"

	"github.com/jackc/pgx/v5"
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

	result := httplib.SuccessfulResult()
	result.Id = credentials.UserId

	httplib.ResponseJSON(w, result)
}

func (s *service) Login(w http.ResponseWriter, r *http.Request) {

}
