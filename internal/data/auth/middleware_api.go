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
			cookie, err := r.Cookie("session_id")
			switch err {
			case nil:
				user, err := userSys.GetUser(cookie.Value)
				switch {
				case err != nil:
					slog.WarnContext(
						r.Context(),
						"Failed to get user from session cookie",
						"channel", "api",
						"error", err.Error(),
					)

				case user != nil:
					serveNext(w, r, user)
					return
				}

			case http.ErrNoCookie:
				slog.WarnContext(
					r.Context(),
					"Session cookie not found",
					"channel", "api",
				)

			default:
				slog.WarnContext(
					r.Context(),
					"Failed to read session cookie",
					"channel", "api",
					"error", err.Error(),
				)
			}

			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer ") {
				serveError(w, r, errors.New("missing token"))
				return
			}

			token := authHeader[len("Bearer "):]

			user, err := tokenSys.VerifyToken(token)
			switch {
			case err != nil:
				serveError(w, r, err)

			case user == nil:
				serveError(w, r, errors.New("invalid token"))

			default:
				serveNext(w, r, user)
			}
		})
	}
}
