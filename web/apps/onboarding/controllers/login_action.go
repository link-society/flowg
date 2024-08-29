package controllers

import (
	"log/slog"
	"time"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"

	"link-society.com/flowg/web/apps/onboarding/templates/views"
)

func LoginAction(
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notifications := []string{}

		err := r.ParseForm()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"Failed to parse form",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Invalid form data")
		} else {
			username := r.FormValue("username")
			password := r.FormValue("password")

			user, err := userSys.GetUser(username)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"Failed to fetch user",
					"channel", "web",
					"error", err.Error(),
				)

				notifications = append(notifications, "&#10060; Invalid credentials")
			} else if user == nil {
				notifications = append(notifications, "&#10060; Invalid credentials")
			} else {
				valid, err := userSys.VerifyUserPassword(username, password)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"Failed to verify user password",
						"channel", "web",
						"error", err.Error(),
					)

					notifications = append(notifications, "&#10060; Invalid credentials")
				} else if !valid {
					notifications = append(notifications, "&#10060; Invalid credentials")
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
			views.LoginProps{
				Notifications: notifications,
			},
		))

		h.ServeHTTP(w, r)
	}
}
