package httplib

import (
	"encoding/json"
	"net/http"
)

func makeError(internalCode int, err error) *ApiError {
	return NewApiError(int(internalCode), err.Error())
}

// ErrorResponseJSON sends error http response
func ErrorResponseJSON(w http.ResponseWriter, httpCode int, internalCode int, err error) {
	writeJSON(w, httpCode, makeError(internalCode, err))
}

// ResponseJSON sends OK JSON response
func ResponseJSON(w http.ResponseWriter, v any) {
	writeJSON(w, http.StatusOK, v)
}

// writeJSON writes a JSON-encoded response to the provided http.ResponseWriter.
// It sets the "Content-Type" header to "application/json; charset=utf-8".
// The HTTP status code is specified by httpCode, and the payload is marshaled from v.
// If an error occurs during marshaling, an error message is written to the response body.
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
