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
	tokenSys := NewTokenSystem(db)
	userSys := NewUserSystem(db)

	return func(next http.Handler) http.Handler {
		serveError := func(w http.ResponseWriter, r *http.Request, err error) {
			code, payload := rest.Err(status.Wrap(err, status.Unauthenticated))
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(payload)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"Failed to encode error response",
					"channel", "api",
					"error", err.Error(),
				)
			}
		}

		serveNext := func(w http.ResponseWriter, r *http.Request, user *User) {
			slog.DebugContext(
				r.Context(),
				"Authenticated user",
				"channel", "api",
				"user", user.Name,
			)

			ctx := ContextWithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			switch {
			case strings.HasPrefix(authHeader, "Basic "):
				token := authHeader[len("Basic "):]

				user, err := tokenSys.VerifyToken(token)
				switch {
				case err != nil:
					serveError(w, r, err)

				case user == nil:
					serveError(w, r, errors.New("invalid token"))

				default:
					serveNext(w, r, user)
				}

			case strings.HasPrefix(authHeader, "Bearer "):
				token := authHeader[len("Bearer "):]

				username, err := VerifyJWT(token)
				if err != nil {
					serveError(w, r, err)
					return
				}

				user, err := userSys.GetUser(username)
				switch {
				case err != nil:
					serveError(w, r, err)

				case user == nil:
					serveError(w, r, errors.New("invalid token"))

				default:
					serveNext(w, r, user)
				}

			default:
				serveError(w, r, errors.New("missing token"))
			}
		})
	}
}
