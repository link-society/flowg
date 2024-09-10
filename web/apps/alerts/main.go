package alerts

import (
	"log/slog"
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"

	"link-society.com/flowg/web/apps/alerts/controllers"
)

func Application(
	authDb *auth.Database,
	configStorage *config.Storage,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)
	alertSys := config.NewAlertSystem(configStorage)

	mux.HandleFunc(
		"GET /web/alerts/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			switch alerts, err := alertSys.List(); {
			case err == nil && len(alerts) > 0:
				url := "/web/alerts/edit/" + alerts[0] + "/"
				http.Redirect(w, r, url, http.StatusSeeOther)

			case err != nil:
				slog.ErrorContext(
					r.Context(),
					"Failed to list alerts",
					"channel", "web",
					"error", err,
				)
				fallthrough

			default:
				http.Redirect(w, r, "/web/alerts/new", http.StatusSeeOther)
			}
		},
	)

	mux.HandleFunc(
		"GET /web/alerts/new/{$}",
		controllers.PageNew(userSys, alertSys),
	)
	mux.HandleFunc(
		"POST /web/alerts/new/{$}",
		controllers.ProcessNewSaveAction(userSys, alertSys),
	)

	mux.HandleFunc(
		"GET /web/alerts/edit/{name}/{$}",
		controllers.PageEdit(userSys, alertSys),
	)
	mux.HandleFunc(
		"POST /web/alerts/edit/{name}/{$}",
		controllers.ProcessEditSaveAction(userSys, alertSys),
	)

	mux.HandleFunc(
		"GET /web/alerts/delete/{name}/{$}",
		controllers.ProcessDeleteAction(userSys, alertSys),
	)

	return mux
}
