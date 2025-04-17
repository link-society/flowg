package auth

import (
	"errors"
	"log/slog"

	"encoding/json"
	"strings"

	"net/http"

	"github.com/swaggest/rest"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/models"

	authUtils "link-society.com/flowg/internal/utils/auth"

	"link-society.com/flowg/internal/storage/auth"
)

func ApiMiddleware(authStorage *auth.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		serveError := func(w http.ResponseWriter, r *http.Request, err error) {
			code, payload := rest.Err(status.Wrap(err, status.Unauthenticated))
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(payload)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"Failed to encode error response",
					slog.String("channel", "api"),
					slog.String("error", err.Error()),
				)
			}
		}

		serveNext := func(w http.ResponseWriter, r *http.Request, user *models.User) {
			slog.DebugContext(
				r.Context(),
				"Authenticated user",
				slog.String("channel", "api"),
				slog.String("user", user.Name),
			)

			ctx := ContextWithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			switch {
			case strings.HasPrefix(strings.ToLower(authHeader), "bearer pat_"):
				token := authHeader[len("bearer "):]

				user, err := authStorage.VerifyToken(r.Context(), token)
				switch {
				case err != nil:
					serveError(w, r, err)

				case user == nil:
					serveError(w, r, errors.New("invalid token"))

				default:
					serveNext(w, r, user)
				}

			case strings.HasPrefix(strings.ToLower(authHeader), "bearer jwt_"):
				token := authHeader[len("bearer "):]

				username, err := authUtils.VerifyJWT(token)
				if err != nil {
					serveError(w, r, err)
					return
				}

				user, err := authStorage.FetchUser(r.Context(), username)
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
