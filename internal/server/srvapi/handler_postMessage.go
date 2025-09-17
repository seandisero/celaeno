package srvapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
	"github.com/seandisero/celaeno/internal/shared"
)

func (api *ApiHandler) HandlerPostMessage(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not read request body", err)
		return
	}

	var params shared.Message
	err = json.Unmarshal(data, &params)
	if err != nil {
		slog.Error("error decoding request body: %v", "error", err)
		server.RespondWithError(w, http.StatusInternalServerError, "could not decode request body", err)
		return
	}

	resp := shared.Message{
		Message: params.Message,
	}

	server.RespondWithJSON(w, http.StatusOK, resp)
}
