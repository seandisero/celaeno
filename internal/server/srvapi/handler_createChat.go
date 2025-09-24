package srvapi

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/seandisero/celaeno/internal/server/chat"
)

func (api *ApiHandler) HandlerCreateChat(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		slog.Error("could not get id from contextv", "error", err)
		return
	}

	user, err := api.DB.GetUserByID(context.Background(), []byte(userID))
	if err != nil {
		slog.Error("could not find user by id")
		return
	}

	slog.Info("creating chat with username", "username", user.Username)

	api.ChatService.Chats[user.Username] = chat.NewChat()
	defer func() {
		delete(api.ChatService.Chats, user.Username)
	}()

	err = api.ChatService.Chats[user.Username].Subscribe(w, r)
	if err != nil {
		if strings.Contains(err.Error(), "context canceled") {
			return
		}
		slog.Error("could not connect to websocket", "error", err)
		return
	}
}
