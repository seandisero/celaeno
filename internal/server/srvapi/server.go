package srvapi

import (
	"net/http"

	"github.com/seandisero/celaeno/internal/server/chat"
	"github.com/seandisero/celaeno/internal/server/database"
)

type ApiHandler struct {
	DB          *database.Queries
	ChatService *chat.ChatService
	JwtSecret   string
}

func (ApiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}
