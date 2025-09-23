package srvapi

import (
	"log/slog"
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
)

func (api *ApiHandler) HandlerStartup(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Status string `json:"status"`
	}{Status: "good to go"}
	slog.Info("started")
	server.RespondWithJSON(w, http.StatusOK, res)
}
