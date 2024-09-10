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

func ProcessEditSaveAction(
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

		_, err := alertSys.Read(r.PathValue("name"))
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch alert", err)
			http.Redirect(w, r, "/web/alerts/new", http.StatusTemporaryRedirect)
			return
		}

		var (
			alertName           string
			webhookUrl          string
			webhookHeaderNames  []string
			webhookHeaderValues []string
			webhook             alerting.Webhook
		)

		if !webutils.Permissions(r.Context()).CanEditAlerts {
			webutils.NotifyError(r.Context(), "You do not have permission to edit alerts")
			goto response
		}

		if err := r.ParseForm(); err != nil {
			webutils.LogError(r.Context(), "Failed to parse form data", err)
			webutils.NotifyError(r.Context(), "Could not parse form data")
			goto response
		}

		alertName = r.FormValue("name")
		webhookUrl = r.FormValue("url")

		webhookHeaderNames = r.Form["header_name"]
		webhookHeaderValues = r.Form["header_value"]

		if alertName == "" {
			webutils.NotifyError(r.Context(), "Alert name is required")
		}

		if webhookUrl == "" {
			webutils.NotifyError(r.Context(), "Webhook URL is required")
		}

		if len(webhookHeaderNames) != len(webhookHeaderValues) {
			webutils.NotifyError(r.Context(), "Header names and values must match")
		}

		if alertName == "" || webhookUrl == "" || len(webhookHeaderNames) != len(webhookHeaderValues) {
			goto response
		}

		webhook.Url = webhookUrl
		webhook.Headers = make(map[string]string, len(webhookHeaderNames))

		for i := range webhookHeaderNames {
			webhook.Headers[webhookHeaderNames[i]] = webhookHeaderValues[i]
		}

		if err := alertSys.Write(alertName, &webhook); err != nil {
			webutils.LogError(r.Context(), "Failed to save alert", err)
			webutils.NotifyError(r.Context(), "Could not save alert")
			goto response
		}

		webutils.NotifyInfo(r.Context(), "Alert saved")

	response:
		alerts, err := alertSys.List()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch alerts", err)
			webutils.NotifyError(r.Context(), "Could not fetch alerts")
			alerts = []string{}
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				Alerts:       alerts,
				CurrentAlert: alertName,
				Webhook:      webhook,
			},
		))
		h.ServeHTTP(w, r)
	}
}
