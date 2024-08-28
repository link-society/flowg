package controllers

import (
	"log/slog"

	"encoding/json"
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/templates/views"
)

func AccountController(authDb *auth.Database) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/account/{$}", func(w http.ResponseWriter, r *http.Request) {
		permissions := auth.Permissions{}
		notifications := []string{}

		username := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(username)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		user, err := authDb.GetUser(username)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error fetching user",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch user details")

			user = &auth.User{
				Name:  username,
				Roles: []string{},
			}
		}

		h := templ.Handler(views.Account(
			views.AccountProps{
				User:       user,
				TokenUUIDs: []string{},
			},
			permissions,
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /web/account/change-password/{$}", func(w http.ResponseWriter, r *http.Request) {
		notifications := []string{}

		username := auth.GetContextUser(r.Context())
		user, err := authDb.GetUser(username)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error fetching user",
				"channel", "web",
				"user", username,
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch user details")
		} else {
			err := r.ParseForm()
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error parsing form",
					"channel", "web",
					"error", err.Error(),
				)

				notifications = append(notifications, "&#10060; Could not parse form")
			} else {
				oldPassword := r.Form.Get("old_password")
				newPassword := r.Form.Get("new_password")

				valid, err := authDb.VerifyUserPassword(user.Name, oldPassword)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"error verifying user password",
						"channel", "web",
						"user", username,
						"error", err.Error(),
					)

					notifications = append(notifications, "&#10060; Could not verify user password")
				}

				if !valid {
					notifications = append(notifications, "&#10060; Invalid password")
				} else {
					err := authDb.SaveUser(
						auth.User{Name: user.Name, Roles: user.Roles},
						newPassword,
					)
					if err != nil {
						slog.ErrorContext(
							r.Context(),
							"error saving user",
							"channel", "web",
							"user", username,
							"error", err.Error(),
						)

						notifications = append(notifications, "&#10060; Could not change user password")
					} else {
						notifications = append(notifications, "&#9989; Password changed")
					}
				}
			}
		}

		trigger := map[string]interface{}{
			"htmx-custom-toast": map[string]interface{}{
				"messages": notifications,
			},
		}

		triggerData, err := json.Marshal(trigger)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error marshalling trigger",
				"channel", "web",
				"error", err.Error(),
			)
		} else {
			w.Header().Add("HX-Trigger", string(triggerData))
		}

		w.WriteHeader(http.StatusOK)
	})

	return mux
}
