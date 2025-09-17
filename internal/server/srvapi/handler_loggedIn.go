package srvapi

import (
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
	"github.com/seandisero/celaeno/internal/shared"
)

func (api *ApiHandler) HandlerLoggedIn(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "could not get user id from context", err)
		return
	}

	user, err := api.DB.GetUserByID(r.Context(), []byte(userID))
	if err != nil {
		server.RespondWithError(w, http.StatusNotFound, "user does not exist", err)
		return
	}

	server.RespondWithJSON(w, http.StatusOK, shared.User{
		ID:          user.ID,
		Username:    user.Username,
		Displayname: user.Displayname,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}
