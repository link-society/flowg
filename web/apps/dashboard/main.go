package dashboard

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/data/pipelines"

	"link-society.com/flowg/web/apps/dashboard/controllers"
)

func Application(
	authDb *auth.Database,
	logDb *logstorage.Storage,
	pipelinesManager *pipelines.Manager,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)
	metaSys := logstorage.NewMetaSystem(logDb)

	mux.HandleFunc(
		"GET /web/{$}",
		controllers.Page(userSys, metaSys, pipelinesManager),
	)

	return mux
}
