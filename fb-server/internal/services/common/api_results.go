package common

import (
	"encoding/json"
	"net/http"

	internalErr "fightbettr.com/fb-server/pkg/errors"
	"fightbettr.com/fb-server/pkg/httplib"
	"fightbettr.com/fb-server/pkg/model"

	"github.com/jackc/pgx/v5"
)

// AddResult handles the addition of a result for a specific fight. It expects a JSON-encoded
// request containing the fight result details. It begins a transaction, sets the fight result,
// checks if the associated event is done, and commits the transaction. If any error occurs during
// the process, it responds with an appropriate API error along with the HTTP status code.
func (s *service) AddResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req model.FightResultRequest
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Events, err)
	}

	tx, err := s.Repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		s.Logger.Errorf("Unable to begin transaction: %s", err)
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Tx, err)
		return
	}

	err = s.Repo.SetFightResult(ctx, tx, &req)
	if err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.EventsFightResult, err)
	}

	err = s.CheckEventIsDone(ctx, tx, req.FightId)
	if err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.EventIsDone, err)
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		s.Logger.Errorf("Unable to commit transaction: %s", txErr)
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.TxCommit, txErr)
		return
	}

	result := httplib.SuccessfulResult()
	result.Id = req.FightId

	httplib.ResponseJSON(w, result)
}
