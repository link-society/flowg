package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"

	"link-society.com/flowg/web/apps/onboarding/templates/views"
)

func LoginForm(
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil && cookie.Value != "" {
			user, err := userSys.GetUser(cookie.Value)
			if err == nil && user != nil {
				http.Redirect(w, r, "/web", http.StatusSeeOther)
				return
			}
		}

		h := templ.Handler(views.Login())

		h.ServeHTTP(w, r)
	}
}
