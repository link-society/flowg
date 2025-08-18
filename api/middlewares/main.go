package middlewares

import (
	"net/http"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
)

type Dependencies struct {
	AuthStorage   auth.Storage
	LogStorage    log.Storage
	ConfigStorage config.Storage

	LogNotifier    lognotify.LogNotifier
	PipelineRunner pipelines.Runner
}

func NewHandler(deps *Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/api/v1/middlewares/elastic/", newElasticHandler(deps))

	return mux
}
