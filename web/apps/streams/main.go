package streams

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/web/apps/streams/controllers"
)

func Application(
	authDb *auth.Database,
	logDb *logstorage.Storage,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)
	metaSys := logstorage.NewMetaSystem(logDb)
	querySys := logstorage.NewQuerySystem(logDb)

	mux.HandleFunc(
		"GET /web/streams/{$}",
		controllers.DefaultPage(userSys, metaSys),
	)
	mux.HandleFunc(
		"GET /web/streams/{name}/{$}",
		controllers.StreamPage(userSys, metaSys, querySys),
	)

	return mux
}
