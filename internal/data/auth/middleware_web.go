package auth

import (
	"log/slog"
	"net/http"

	"link-society.com/flowg/internal/webutils/htmx"
)

func WebMiddleware(db *Database) func(http.Handler) http.Handler {
	userSys := NewUserSystem(db)

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

				if htmx.IsHtmxRequest(r) {
					w.Header().Set("HX-Redirect", "/auth/login")
					w.WriteHeader(http.StatusOK)
				} else {
					http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				}
				return
			}

			user, err := userSys.GetUser(cookie.Value)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"Failed to get user from session cookie",
					"channel", "web",
					"error", err.Error(),
				)

				if htmx.IsHtmxRequest(r) {
					w.Header().Set("HX-Redirect", "/auth/login")
					w.WriteHeader(http.StatusOK)
				} else {
					http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				}
				return
			}

			if user == nil {
				slog.ErrorContext(
					r.Context(),
					"User not found",
					"channel", "web",
				)

				if htmx.IsHtmxRequest(r) {
					w.Header().Set("HX-Redirect", "/auth/login")
					w.WriteHeader(http.StatusOK)
				} else {
					http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				}
				return
			}

			slog.DebugContext(
				r.Context(),
				"Authenticated user",
				"channel", "web",
				"user", user.Name,
			)

			ctx := ContextWithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
