package streams

import (
	"net/http"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/logstorage"

	"link-society.com/flowg/web/apps/streams/controllers"
)

func Application(
	authDb *auth.Database,
	logDb *logstorage.Storage,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)

	mux.HandleFunc(
		"GET /web/streams/{$}",
		controllers.Index(userSys, logDb),
	)
	mux.HandleFunc(
		"GET /web/streams/{name}/{$}",
		controllers.Detail(userSys, logDb),
	)

	return mux
}
