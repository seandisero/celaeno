package srvapi

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
)

func (api *ApiHandler) HandlerPostMessage(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		slog.Error("could not get path value", "name", name)
		server.RespondWithError(w, http.StatusInternalServerError, "could not get user name from request", nil)
		return
	}

	body := http.MaxBytesReader(w, r.Body, 8192)
	data, err := io.ReadAll(body)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not read request body", err)
		return
	}

	chat, ok := api.ChatService.Chats[name]
	if !ok {
		slog.Error("chat does not exist", "username", name)
		server.RespondWithError(w, http.StatusBadRequest, "chat does not exist for user", err)
		return
	}
	chat.PublishMessage(data)

	server.RespondWithJSON(w, 200, http.NoBody)
}
