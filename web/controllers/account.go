package controllers

import (
	"log/slog"

	"encoding/json"
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/templates/components"
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
				"user", username,
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch user details")

			user = &auth.User{
				Name:  username,
				Roles: []string{},
			}
		}

		tokenUUIDs, err := authDb.ListPersonalAccessTokens(user.Name)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing personal access tokens",
				"channel", "web",
				"user", username,
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch personal access tokens")
			tokenUUIDs = []string{}
		}

		h := templ.Handler(views.Account(
			views.AccountProps{
				User:       user,
				TokenUUIDs: tokenUUIDs,
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

	mux.HandleFunc("POST /web/account/token/new/{$}", func(w http.ResponseWriter, r *http.Request) {
		success := false
		notifications := []string{}
		username := auth.GetContextUser(r.Context())

		token, err := auth.NewToken(32)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error generating token",
				"channel", "web",
				"user", username,
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not generate token")
		} else {
			err := authDb.AddPersonalAccessToken(username, token)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error saving token",
					"channel", "web",
					"user", username,
					"error", err.Error(),
				)

				notifications = append(notifications, "&#10060; Could not save token")
			} else {
				success = true
			}
		}

		trigger := map[string]interface{}{
			"htmx-custom-toast": map[string]interface{}{
				"messages": notifications,
			},
		}

		if success {
			trigger["htmx-custom-modal-open"] = map[string]interface{}{}
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

		h := templ.Handler(components.TokenViewer(token))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /web/account/token/delete/{tokenUUID}/{$}", func(w http.ResponseWriter, r *http.Request) {
		username := auth.GetContextUser(r.Context())
		notifications := []string{}

		tokenUUID := r.PathValue("tokenUUID")
		err := authDb.DeletePersonalAccessToken(username, tokenUUID)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error deleting token",
				"channel", "web",
				"user", username,
				"tokenUUID", tokenUUID,
				"error", err.Error(),
			)

			w.Header().Add("HX-Reswap", "none")

			notifications = append(notifications, "&#10060; Could not delete token")
		} else {
			w.Header().Add("HX-Reswap", "delete")
			w.Header().Add("HX-Retarget", "tr[data-token='"+tokenUUID+"']")

			notifications = append(notifications, "&#9989; Token deleted")
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
