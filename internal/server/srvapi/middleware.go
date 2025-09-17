package srvapi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/seandisero/celaeno/internal/server"
	"github.com/seandisero/celaeno/internal/server/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

func (api ApiHandler) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			server.RespondWithError(w, http.StatusUnauthorized, "no authorization token", err)
			return
		}

		userID, err := auth.ValidateJWT(token, api.JwtSecret)
		if err != nil {
			slog.Error("could not validate jwt")
			server.RespondWithError(w, http.StatusUnauthorized, "invalid token", err)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	slog.Info("getting id from context", "context", ctx)
	userID, ok := ctx.Value(UserContextKey).(string)
	if !ok {
		return "", fmt.Errorf("UserContextKey not present in context")
	}

	return userID, nil
}
