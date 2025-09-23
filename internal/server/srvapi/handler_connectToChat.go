package srvapi

import (
	"log/slog"
	"net/http"
)

func (api *ApiHandler) HandlerConnectToChat(w http.ResponseWriter, r *http.Request) {
	connectionName := r.PathValue("name")

	chat, ok := api.ChatService.Chats[connectionName]
	if !ok {
		slog.Error("user is attempting to connect with chat that does not exist", "id", connectionName)
		return
	}

	err := chat.Subscribe(w, r)
	if err != nil {
		slog.Error("user no longer subscribed", "error", err)
		return
	}
}
