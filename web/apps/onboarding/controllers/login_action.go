package controllers

import (
	"time"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"

	"link-society.com/flowg/web/apps/onboarding/templates/views"
)

func LoginAction(
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))

		var (
			username string
			password string
		)

		if err := r.ParseForm(); err != nil {
			webutils.LogError(r.Context(), "Failed to parse form", err)
			webutils.NotifyError(r.Context(), "Invalid form data")
			goto response
		}

		username = r.FormValue("username")
		password = r.FormValue("password")

		switch user, err := userSys.GetUser(username); {
		case err != nil:
			webutils.LogError(r.Context(), "Failed to fetch user", err)
			webutils.NotifyError(r.Context(), "Invalid credentials")
			goto response

		case user == nil:
			webutils.NotifyError(r.Context(), "Invalid credentials")
			goto response
		}

		switch valid, err := userSys.VerifyUserPassword(username, password); {
		case err != nil:
			webutils.LogError(r.Context(), "Failed to verify user password", err)
			webutils.NotifyError(r.Context(), "Invalid credentials")
			goto response

		case !valid:
			webutils.NotifyError(r.Context(), "Invalid credentials")
			goto response

		case valid:
			cookie := &http.Cookie{
				Name:     "session_id",
				Value:    username,
				Expires:  time.Now().Add(24 * time.Hour),
				Path:     "/",
				HttpOnly: true,
			}

			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

	response:
		h := templ.Handler(views.Login())

		h.ServeHTTP(w, r)
	}
}
