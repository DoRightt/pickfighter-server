package common

import (
	"encoding/json"
	"net/http"
	internalErr "projects/fb-server/errors"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"

	"github.com/jackc/pgx/v5"
)

func (s *service) AddResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req model.FightResultRequest
	if err := decoder.Decode(&req); err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Events, err)
	}

	tx, err := s.Repo.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		s.Logger.Errorf("Unable to begin transaction: %s", err)
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Tx, err)
	}

	err = s.Repo.SetFightResult(ctx, tx, &req)
	if err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.EventsFightResult, err)
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
