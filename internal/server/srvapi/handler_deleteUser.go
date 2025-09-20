package srvapi

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/seandisero/celaeno/internal/server"
	"github.com/seandisero/celaeno/internal/server/auth"
)

func (api *ApiHandler) HandlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var req request
	err := decoder.Decode(&req)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not decode body", err)
		return
	}

	// NOTE: I realized that if I want to be able to delete other users as an admin later without goin onto turso
	// then I cant just use the auth token id to delet the user.
	userID, err := GetUserIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "invalid uuid", err)
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "invalid user uuid", err)
		return
	}

	if userUUID != id {
		server.RespondWithError(w, http.StatusInternalServerError, "user id does not match id of delete request", err)
		return
	}

	// NOTE: this will need to be adjusted later when an admin wants to delete a user, it shouldn't require a password
	user, err := api.DB.GetUserByID(r.Context(), []byte(id.String()))
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not get user by id", err)
		return
	}

	err = auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "passwrod does not match", err)
		return
	}

	err = api.DB.DeleteUserByID(r.Context(), []byte(id.String()))
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "error deleting user", err)
		return
	}

	server.RespondWithJSON(w, http.StatusNoContent, http.NoBody)
}
