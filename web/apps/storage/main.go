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
		"GET /web/storage/edit/{name}/{$}",
		controllers.StreamPage(userSys, metaSys),
	)
	mux.Handle(
		"POST /web/storage/delete/{name}/{$}",
		controllers.ProcessStreamDeleteAction(userSys, metaSys),
	)

	mux.Handle(
		"GET /web/storage/new/{$}",
		controllers.DisplayStreamCreateForm(userSys),
	)
	mux.Handle(
		"POST /web/storage/new/{$}",
		controllers.ProcessStreamCreateForm(userSys, metaSys),
	)

	mux.Handle(
		"POST /web/storage/edit/{name}/retention/{$}",
		controllers.ProcessRetentionSaveAction(userSys, metaSys),
	)
	mux.Handle(
		"GET /web/storage/edit/{name}/index/add/{$}",
		controllers.ProcessIndexAddAction(userSys),
	)
	mux.Handle(
		"POST /web/storage/edit/{name}/index/add/{$}",
		controllers.ProcessIndexAddForm(userSys, metaSys),
	)
	mux.Handle(
		"GET /web/storage/edit/{name}/index/delete/{field}/{$}",
		controllers.ProcessIndexDeleteAction(userSys, metaSys),
	)

	return mux
}
