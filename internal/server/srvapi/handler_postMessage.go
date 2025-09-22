package srvapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
	"github.com/seandisero/celaeno/internal/shared"
)

func (api *ApiHandler) HandlerPostMessage(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 8192)
	data, err := io.ReadAll(body)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not read request body", err)
		return
	}

	var msg shared.Message
	err = json.Unmarshal(data, &msg)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not unmarshal request body", err)
		return
	}

	for _, chat := range api.ChatService.Chats {
		chat.PublishMessage(data)
	}

	server.RespondWithJSON(w, 200, http.NoBody)
}
