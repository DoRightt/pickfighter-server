package httplib

import (
	"encoding/json"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, v any) {
	writeJSON(w, http.StatusOK, v)
}

func writeJSON(w http.ResponseWriter, httpCode int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response, err := json.Marshal(v)
	if err != nil {
		w.Write([]byte("Cannot marshal response: " + err.Error()))
		return
	}

	w.WriteHeader(httpCode)
	w.Write(response)
}
