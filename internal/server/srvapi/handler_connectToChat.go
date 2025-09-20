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
		server.RespondWithError(w, http.StatusInternalServerError, "chat does not seem to exist", nil)

	}

	chat.Subscribe(w, r)

	server.RespondWithJSON(w, http.StatusOK, http.NoBody)
}
