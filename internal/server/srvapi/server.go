package srvapi

import (
	"net/http"

	"github.com/seandisero/celaeno/internal/server/database"
)

type ApiHandler struct {
	DB        *database.Queries
	JwtSecret string
}

func (ApiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}
