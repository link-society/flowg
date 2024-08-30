package controllers

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"
)

func ChangePassword(
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		user := auth.GetContextUser(r.Context())

		var (
			oldPassword string
			newPassword string
		)

		if err := r.ParseForm(); err != nil {
			webutils.LogError(r.Context(), "Failed to parse form data", err)
			webutils.NotifyError(r.Context(), "Could not parse form")
			goto response
		}

		oldPassword = r.Form.Get("old_password")
		newPassword = r.Form.Get("new_password")

		switch valid, err := userSys.VerifyUserPassword(user.Name, oldPassword); {
		case err != nil:
			webutils.LogError(r.Context(), "Failed to verify user password", err)
			webutils.NotifyError(r.Context(), "Could not verify user password")
			goto response

		case !valid:
			webutils.NotifyError(r.Context(), "Invalid password")
			goto response
		}

		if err := userSys.SaveUser(*user, newPassword); err != nil {
			webutils.LogError(r.Context(), "Failed to save user", err)
			webutils.NotifyError(r.Context(), "Could not change user password")
			goto response
		}

		webutils.NotifyInfo(r.Context(), "Password changed")

	response:
		trigger := htmx.Trigger{
			ToastEvent: &htmx.ToastEvent{
				Messages: webutils.Notifications(r.Context()),
			},
		}

		trigger.Write(r.Context(), w)
		w.WriteHeader(http.StatusOK)
	}
}
