package srvapi

import (
	"log/slog"
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
)

func (api ApiHandler) HandlerStatus(w http.ResponseWriter, r *http.Request) {
	for key := range api.ChatService.Chats {
		slog.Info("chat name", "id", key)
	}
	server.RespondWithJSON(w, 200, http.NoBody)
}
