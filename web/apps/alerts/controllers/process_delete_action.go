package controllers

import (
	"fmt"

	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/webutils"
)

func ProcessDeleteAction(
	userSys *auth.UserSystem,
	alertSys *config.AlertSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))
		alertName := r.PathValue("name")

		if !webutils.Permissions(r.Context()).CanViewAlerts {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		} else if !webutils.Permissions(r.Context()).CanEditAlerts {
			url := fmt.Sprintf("/web/alerts/edit/%s", alertName)
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}

		if err := alertSys.Delete(alertName); err != nil {
			webutils.LogError(r.Context(), "Failed to delete alert", err)
		}

		http.Redirect(w, r, "/web/alerts/new", http.StatusTemporaryRedirect)
	}
}
