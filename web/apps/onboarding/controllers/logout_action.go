package controllers

import (
	"time"

	"net/http"
)

func LogoutAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
