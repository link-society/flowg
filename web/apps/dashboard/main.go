package dashboard

import (
	"net/http"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/logstorage"
	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/apps/dashboard/controllers"
)

func Application(
	authDb *auth.Database,
	logDb *logstorage.Storage,
	pipelinesManager *pipelines.Manager,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)

	mux.HandleFunc(
		"GET /web/{$}",
		controllers.Index(userSys, logDb, pipelinesManager),
	)

	return mux
}
