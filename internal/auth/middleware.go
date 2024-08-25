package auth

import (
	"errors"
	"log/slog"

	"encoding/json"
	"strings"

	"net/http"

	"github.com/swaggest/rest"
	"github.com/swaggest/usecase/status"
)

func ApiMiddleware(db *Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer ") {
				code, payload := rest.Err(status.Wrap(
					errors.New("missing token"),
					status.Unauthenticated,
				))
				w.WriteHeader(code)
				err := json.NewEncoder(w).Encode(payload)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"Failed to encode error response",
						"channel", "api",
						"error", err.Error(),
					)
				}
				return
			}

			token := authHeader[len("Bearer "):]

			username, found, err := db.VerifyPersonalAccessToken(token)
			if err != nil {
				code, payload := rest.Err(status.Wrap(err, status.Unauthenticated))
				w.WriteHeader(code)
				err := json.NewEncoder(w).Encode(payload)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"Failed to encode error response",
						"channel", "api",
						"error", err.Error(),
					)
				}
				return
			}

			if !found {
				code, payload := rest.Err(status.Wrap(
					errors.New("invalid token"),
					status.Unauthenticated,
				))
				w.WriteHeader(code)
				err := json.NewEncoder(w).Encode(payload)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"Failed to encode error response",
						"channel", "api",
						"error", err.Error(),
					)
				}
				return
			}

			ctx := ContextWithUsername(r.Context(), username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
