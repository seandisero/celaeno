package srvapi

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
)

func (api *ApiHandler) HandlerConnectToChat(w http.ResponseWriter, r *http.Request) {
	connectionName := r.PathValue("name")

	slog.Info("connectionName", "name", connectionName)
	user, err := api.DB.GetUserByName(context.Background(), connectionName)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not Get user by name", err)
		return
	}

	userID := string(user.ID)

	slog.Info("making user id into string", "id", userID)

	chat, ok := api.ChatService.Chats[userID]
	if !ok {
		slog.Error("user is attempting to connect with chat that does not exist", "id", userID)
		return
	}

	err = chat.Subscribe(w, r)
	if err != nil {
		slog.Error("user fuled to subscribe")
		return
	}
}
