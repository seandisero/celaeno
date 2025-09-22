package srvapi

import (
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

	newChat := chat.NewChat()
	api.ChatService.Chats[userID] = newChat

	err = api.ChatService.Chats[userID].Subscribe(w, r)
	if err != nil {
		server.RespondWithError(w, 400, "could not connect to websocket", err)
		return
	}
}
