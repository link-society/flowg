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

			slog.DebugContext(
				r.Context(),
				"Authenticated user",
				"channel", "api",
				"user", username,
			)

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
		user := ctx.Value(CONTEXT_USERNAME).(string)
		authorized, err := db.VerifyUserPermission(
			user,
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

		slog.DebugContext(
			ctx,
			"Authorized user",
			"channel", "api",
			"user", user,
			"scope", scope,
		)

		return next(ctx, req, resp)
	}
}

func WebMiddleware(db *Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				if err == http.ErrNoCookie {
					slog.ErrorContext(
						r.Context(),
						"Session cookie not found",
						"channel", "web",
					)
				} else {
					slog.ErrorContext(
						r.Context(),
						"Failed to read session cookie",
						"channel", "web",
						"error", err.Error(),
					)
				}

				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			user, err := db.GetUser(cookie.Value)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"Failed to get user from session cookie",
					"channel", "web",
					"error", err.Error(),
				)

				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			if user == nil {
				slog.ErrorContext(
					r.Context(),
					"User not found",
					"channel", "web",
				)

				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			slog.DebugContext(
				r.Context(),
				"Authenticated user",
				"channel", "web",
				"user", user.Name,
			)

			ctx := ContextWithUsername(r.Context(), user.Name)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
