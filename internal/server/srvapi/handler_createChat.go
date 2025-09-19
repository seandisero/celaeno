package srvapi

import (
	"log/slog"
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
	"github.com/seandisero/celaeno/internal/server/chat"
)

func (api *ApiHandler) HandlerCreateChat(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())

	if err != nil {
		server.RespondWithError(w, 400, "could not connect to websocket", err)
		return
	}

	slog.Info("made connection to websocket")

	newChat := chat.NewChat()
	api.ChatService.Chats[userID] = newChat

	api.ChatService.Chats[userID].Subscribe(w, r)
}
