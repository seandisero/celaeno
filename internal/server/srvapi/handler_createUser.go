package srvapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/seandisero/celaeno/internal/server"
	"github.com/seandisero/celaeno/internal/server/auth"
	"github.com/seandisero/celaeno/internal/server/database"
	"github.com/seandisero/celaeno/internal/shared"
)

func (api *ApiHandler) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var req request
	err := decoder.Decode(&req)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not decode request", err)
		return
	}

	_, err = api.DB.GetUserByName(r.Context(), req.Name)
	if err == nil {
		slog.Error("user already exists")
		server.RespondWithError(w, http.StatusConflict, "user already exists", nil)
		return
	}

	id := uuid.New()
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		slog.Error("trying to hash passwrod", "error", err)
		server.RespondWithError(w, http.StatusInternalServerError, "server error", err)
		return
	}

	userParams := database.CreateUserParams{
		ID:             []byte(id.String()),
		Username:       req.Name,
		HashedPassword: hashedPassword,
	}

	user, err := api.DB.CreateUser(context.Background(), userParams)
	if err != nil {
		slog.Error("failed to create user")
		server.RespondWithError(w, http.StatusInternalServerError, "could not create user", nil)
		return
	}

	server.RespondWithJSON(w, http.StatusOK, shared.User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
