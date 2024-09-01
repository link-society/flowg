package storage

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/web/apps/storage/controllers"
)

func Application(
	authDb *auth.Database,
	logDb *logstorage.Storage,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)
	metaSys := logstorage.NewMetaSystem(logDb)

	mux.Handle(
		"GET /web/storage/{$}",
		controllers.DefaultPage(userSys, metaSys),
	)
	mux.Handle(
		"GET /web/storage/{name}/{$}",
		controllers.StreamPage(userSys, metaSys),
	)

	mux.Handle(
		"POST /web/storage/{name}/retention/{$}",
		controllers.ProcessRetentionSaveAction(userSys, metaSys),
	)

	return mux
}
