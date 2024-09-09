package dashboard

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/web/apps/dashboard/controllers"
)

func Application(
	authDb *auth.Database,
	logDb *logstorage.Storage,
	configStorage *config.Storage,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)
	metaSys := logstorage.NewMetaSystem(logDb)
	transformerSys := config.NewTransformerSystem(configStorage)
	pipelineSys := config.NewPipelineSystem(configStorage)

	mux.HandleFunc(
		"GET /web/{$}",
		controllers.Page(userSys, metaSys, transformerSys, pipelineSys),
	)

	return mux
}
