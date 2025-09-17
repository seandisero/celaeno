package srvapi

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/seandisero/celaeno/internal/server"
	"github.com/seandisero/celaeno/internal/server/database"
	"github.com/seandisero/celaeno/internal/shared"
)

func (api *ApiHandler) HandlerSetDisplayName(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "invalid uuid", err)
		return
	}

	type request struct {
		DisplayName string `json:"displayname"`
	}

	var requestBody request
	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "could not decode request body", err)
		return
	}

	setParams := database.SetUserDisplayNameParams{
		ID: []byte(id.String()),
		Displayname: sql.NullString{
			String: requestBody.DisplayName,
			Valid:  true,
		},
	}

	dbUser, err := api.DB.SetUserDisplayName(r.Context(), setParams)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not set user diaplay name", err)
		return
	}

	user := shared.User{
		ID:          dbUser.ID,
		Username:    dbUser.Username,
		Displayname: dbUser.Displayname,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
	}

	server.RespondWithJSON(w, http.StatusCreated, user)

}
