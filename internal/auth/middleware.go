package auth

import (
	"context"

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

func RequireScopeApiMiddleware[Req any, Resp any](
	db *Database,
	scope Scope,
	next func(context.Context, Req, *Resp) error,
) func(context.Context, Req, *Resp) error {
	return func(ctx context.Context, req Req, resp *Resp) error {
		authorized, err := db.VerifyUserPermission(
			ctx.Value(CONTEXT_USERNAME).(string),
			scope,
		)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"Failed to verify user permission",
				"channel", "api",
				"error", err.Error(),
			)
			return status.Wrap(err, status.PermissionDenied)
		}

		if !authorized {
			return status.PermissionDenied
		}

		return next(ctx, req, resp)
	}
}