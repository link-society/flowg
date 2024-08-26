package controllers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/templates/views"
)

func AuthController(authDb *auth.Database) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /auth/login/{$}", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil && cookie.Value != "" {
			user, err := authDb.GetUser(cookie.Value)
			if err == nil && user != nil {
				http.Redirect(w, r, "/web", http.StatusSeeOther)
				return
			}
		}

		h := templ.Handler(views.Login(
			views.LoginProps{},
			[]string{},
		))

		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /auth/login/{$}", func(w http.ResponseWriter, r *http.Request) {
		notifications := []string{}

		err := r.ParseForm()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"Failed to parse form",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Invalid form data")
		} else {
			username := r.FormValue("username")
			password := r.FormValue("password")

			user, err := authDb.GetUser(username)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"Failed to fetch user",
					"channel", "web",
					"error", err.Error(),
				)

				notifications = append(notifications, "❌ Invalid credentials")
			} else if user == nil {
				notifications = append(notifications, "❌ Invalid credentials")
			} else {
				valid, err := authDb.VerifyUserPassword(username, password)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"Failed to verify user password",
						"channel", "web",
						"error", err.Error(),
					)

					notifications = append(notifications, "❌ Invalid credentials")
				} else if !valid {
					notifications = append(notifications, "❌ Invalid credentials")
				} else {
					cookie := &http.Cookie{
						Name:     "session_id",
						Value:    user.Name,
						Expires:  time.Now().Add(24 * time.Hour),
						Path:     "/",
						HttpOnly: true,
					}

					http.SetCookie(w, cookie)
					http.Redirect(w, r, "/web", http.StatusSeeOther)
					return
				}
			}
		}

		h := templ.Handler(views.Login(
			views.LoginProps{},
			notifications,
		))

		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /auth/logout/{$}", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			Path:     "/",
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	})

	return mux
}
