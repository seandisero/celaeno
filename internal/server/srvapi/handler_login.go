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

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginRequest)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "could not decode request", err)
		return
	}
	defer r.Body.Close()

	user, err := api.DB.GetUserByName(r.Context(), loginRequest.Name)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "could not find user", err)
		return
	}

	err = auth.CheckPasswordHash(loginRequest.Password, user.HashedPassword)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "password is incorrect", nil)
		return
	}

	token, err := auth.MakeJWT(user.ID, api.JwtSecret, 30*time.Second)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not sign string for web token", err)
		return
	}

	type loginResponce struct {
		User  shared.User
		Token string `json:"token"`
	}

	server.RespondWithJSON(w, http.StatusOK, loginResponce{
		User: shared.User{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	})
}
