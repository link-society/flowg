package controllers

import (
	"log/slog"

	"encoding/json"

	"net/http"

	"link-society.com/flowg/internal/auth"
)

func ChangePassword(
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
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

			valid, err := userSys.VerifyUserPassword(user.Name, oldPassword)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error verifying user password",
					"channel", "web",
					"user", user.Name,
					"error", err.Error(),
				)

				notifications = append(notifications, "&#10060; Could not verify user password")
			}

			if !valid {
				notifications = append(notifications, "&#10060; Invalid password")
			} else {
				err := userSys.SaveUser(
					*user,
					newPassword,
				)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"error saving user",
						"channel", "web",
						"user", user.Name,
						"error", err.Error(),
					)

					notifications = append(notifications, "&#10060; Could not change user password")
				} else {
					notifications = append(notifications, "&#9989; Password changed")
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
	}
}
