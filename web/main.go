package web

import (
	"embed"

	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/web/apps/account"
	"link-society.com/flowg/web/apps/admin"
	"link-society.com/flowg/web/apps/alerts"
	"link-society.com/flowg/web/apps/dashboard"
	"link-society.com/flowg/web/apps/onboarding"
	"link-society.com/flowg/web/apps/pipelines"
	"link-society.com/flowg/web/apps/storage"
	"link-society.com/flowg/web/apps/streams"
	"link-society.com/flowg/web/apps/transformers"
)

//go:embed static/**/*.css
//go:embed static/**/*.woff2
//go:embed static/**/*.js
var staticfiles embed.FS

//go:generate templ generate

func NewHandler(
	authDb *auth.Database,
	logDb *logstorage.Storage,
	configStorage *config.Storage,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServer(http.FS(staticfiles)))

	mux.Handle("/auth/", onboarding.Application(authDb))

	authMiddleware := auth.WebMiddleware(authDb)

	mux.Handle(
		"/web/",
		authMiddleware(dashboard.Application(authDb, logDb, configStorage)),
	)
	mux.Handle(
		"/web/streams/",
		authMiddleware(streams.Application(authDb, logDb)),
	)
	mux.Handle(
		"/web/transformers/",
		authMiddleware(transformers.Application(authDb, configStorage)),
	)
	mux.Handle(
		"/web/pipelines/",
		authMiddleware(pipelines.Application(authDb, configStorage)),
	)
	mux.Handle(
		"/web/alerts/",
		authMiddleware(alerts.Application(authDb, configStorage)),
	)
	mux.Handle(
		"/web/storage/",
		authMiddleware(storage.Application(authDb, logDb)),
	)
	mux.Handle(
		"/web/admin/",
		authMiddleware(admin.Application(authDb)),
	)
	mux.Handle(
		"/web/account/",
		authMiddleware(account.Application(authDb)),
	)

	return mux
}
