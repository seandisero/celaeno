package srvapi

import (
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
)

func (api ApiHandler) HandlerListChatRooms(w http.ResponseWriter, r *http.Request) {
	type responceBody struct {
		ChatRooms []string `json:"chat_rooms"`
	}

	roomList := make([]string, 0)
	for room := range api.ChatService.Chats {
		roomList = append(roomList, room)
	}

	if len(roomList) > 20 {
		roomList = roomList[:20]
	}

	respBody := responceBody{
		ChatRooms: roomList,
	}

	server.RespondWithJSON(w, http.StatusOK, respBody)
}
