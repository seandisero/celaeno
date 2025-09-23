package srvapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/seandisero/celaeno/internal/server"
	"github.com/seandisero/celaeno/internal/server/auth"
	"github.com/seandisero/celaeno/internal/shared"
)

func (api *ApiHandler) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest shared.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "could not decode request", err)
		return
	}

	user, err := api.DB.GetUserByName(r.Context(), loginRequest.Name)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "authentication error", err)
		return
	}

	err = auth.CheckPasswordHash(loginRequest.Password, user.HashedPassword)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "password is incorrect", nil)
		return
	}

	token, err := auth.MakeJWT(user.ID, api.JwtSecret, 8*time.Hour)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not make json web token", err)
		return
	}

	type loginResponce struct {
		User  shared.User
		Token string `json:"token"`
	}

	server.RespondWithJSON(w, http.StatusOK, loginResponce{
		User: shared.User{
			ID:          user.ID,
			Username:    user.Username,
			Displayname: user.Displayname,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
		Token: token,
	})
}
