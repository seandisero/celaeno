package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/seandisero/celaeno/internal/shared"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("error while marshaling json", "error", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		slog.Error("could not write data to response", "error", err)
		w.WriteHeader(400)
		return
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string, logErr error) {
	if logErr != nil {
		slog.Error(logErr.Error())
	}
	resp := shared.ResponceError{
		Error: message,
	}
	RespondWithJSON(w, code, resp)
}
