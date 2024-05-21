package http

import (
	"log"
	"net/http"

	"fightbettr.com/pkg/httplib"
	"fightbettr.com/pkg/utils"
)

// GetFighters handles HTTP requests to retrieve fighters based on status.
func (h *Handler) GetFighters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := utils.Capitalize(r.FormValue("status"))

	fighters, err := h.ctrl.SearchFighters(ctx, status)
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
