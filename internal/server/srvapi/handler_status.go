package srvapi

import (
	"log/slog"
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
)

func (api ApiHandler) HandlerStatus(w http.ResponseWriter, r *http.Request) {
	slog.Info("getting status")
	slog.Info("Chats", "ApiHandler", api.ChatService.Chats)
	for name, chat := range api.ChatService.Chats {
		slog.Info("Chat", "name", name)
		slog.Info("Chat", "chat", chat)
	}
	server.RespondWithJSON(w, 200, http.NoBody)
}
