package common

import (
	"encoding/json"
	"net/http"
	internalErr "projects/fb-server/pkg/errors"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"

	"github.com/jackc/pgx/v5"
)

// HandleNewEvent handles the creation of a new event. It decodes the JSON request body into the
// model.EventsRequest struct, begins a new transaction, and calls the CreateEvent method to create
// the event and its associated fights. If any error occurs during the process, it responds with
// an appropriate API error along with the HTTP status code. If the transaction is successful, it
// responds with a successful result and the ID of the created event.
func (s *service) HandleNewEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var req model.EventsRequest
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

	event, err := s.CreateEvent(ctx, tx, &req)
	if err != nil {
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.Events, err)
	}

	if txErr := tx.Commit(ctx); txErr != nil {
		s.Logger.Errorf("Unable to commit transaction: %s", txErr)
		httplib.ErrorResponseJSON(w, http.StatusBadRequest, internalErr.TxCommit, txErr)
		return
	}

	result := httplib.SuccessfulResult()
	result.Id = event.EventId

	httplib.ResponseJSON(w, result)
}

// GetEvents retrieves a list of events. It queries the repository for the count of events and the
// list of events. If any error occurs during the process, it responds with an appropriate API error
// along with the HTTP status code. If no events are found, it responds with an empty list result.
func (s *service) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	count, err := s.Repo.SearchEventsCount(ctx)
	if err != nil {
		s.Logger.Errorf("Failed to get events count: %s", err)
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.CountEvents, err)
		return
	}
	if count == 0 {
		httplib.ResponseJSON(w, httplib.ListResult{})
		return
	}

	events, err := s.Repo.SearchEvents(ctx)
	if err != nil {
		s.Logger.Errorf("Failed to find events: %s", err)
		httplib.ErrorResponseJSON(w, http.StatusInternalServerError, internalErr.Events, err)
		return
	}

	httplib.ResponseJSON(w, httplib.ListResult{
		Results: events,
		Count:   count,
	})
}
