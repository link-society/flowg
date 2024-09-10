package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/alerting"
	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/webutils"

	"link-society.com/flowg/web/apps/alerts/templates/views"
)

func PageNew(
	userSys *auth.UserSystem,
	alertSys *config.AlertSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanViewAlerts {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		alerts, err := alertSys.List()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch alerts", err)
			webutils.NotifyError(r.Context(), "Could not fetch alerts")
			alerts = []string{}
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				Alerts:       alerts,
				CurrentAlert: "",

				Webhook: alerting.Webhook{
					Url:     "",
					Headers: map[string]string{},
				},
			},
		))
		h.ServeHTTP(w, r)
	}
}
