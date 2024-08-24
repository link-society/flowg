package web

import (
	"embed"

	"net/http"

	"link-society.com/flowg/internal/logstorage"
	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/controllers"
)

//go:embed static/**/*.css
//go:embed static/**/*.woff2
//go:embed static/**/*.js
var staticfiles embed.FS

//go:generate templ generate

func NewHandler(
	logDb *logstorage.Storage,
	pipelinesManager *pipelines.Manager,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServer(http.FS(staticfiles)))

	mux.Handle("/web/", controllers.MainController(logDb, pipelinesManager))
	mux.Handle("/web/streams/", controllers.StreamController(logDb))
	mux.Handle("/web/transformers/", controllers.TransformersController(pipelinesManager))
	mux.Handle("/web/pipelines/", controllers.PipelinesController(pipelinesManager))

	return mux
}
