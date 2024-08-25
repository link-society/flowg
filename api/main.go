package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	"github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"

	"link-society.com/flowg/internal"
	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/logstorage"
	"link-society.com/flowg/internal/pipelines"
)

func NewHandler(
	authDb *auth.Database,
	logDb *logstorage.Storage,
	pipelinesManager *pipelines.Manager,
) http.Handler {
	service := web.NewService(openapi31.NewReflector())

	service.OpenAPISchema().SetTitle("Flowg API")
	service.OpenAPISchema().SetVersion(internal.FLOWG_VERSION)
	service.Docs("/api/docs", v5emb.New)

	service.With(
		nethttp.HTTPBearerSecurityMiddleware(
			service.OpenAPICollector,
			"tokenAuth",
			"Authentication using Personal Access Token",
			"PAT",
		),
		auth.ApiMiddleware(authDb),
	).Group(func(router chi.Router) {
		r := &routerWrapper{Router: router}

		r.Get("/api/v1/transformers", ListTransformersUsecase(pipelinesManager))
		r.Get("/api/v1/transformers/{transformer}", GetTransformerUsecase(pipelinesManager))
		r.Put("/api/v1/transformers/{transformer}", SaveTransformerUsecase(pipelinesManager))
		r.Delete("/api/v1/transformers/{transformer}", DeleteTransformerUsecase(pipelinesManager))
		r.Post("/api/v1/transformers/{transformer}/test", TestTransformerUsecase(pipelinesManager))

		r.Get("/api/v1/pipelines", ListPipelinesUsecase(pipelinesManager))
		r.Get("/api/v1/pipelines/{pipeline}", GetPipelineUsecase(pipelinesManager))
		r.Put("/api/v1/pipelines/{pipeline}", SavePipelineUsecase(pipelinesManager))
		r.Delete("/api/v1/pipelines/{pipeline}", DeletePipelineUsecase(pipelinesManager))
		r.Post("/api/v1/pipelines/{pipeline}/logs", IngestLogUsecase(pipelinesManager))

		r.Get("/api/v1/streams", ListStreamsUsecase(logDb))
		r.Get("/api/v1/streams/{stream}", QueryStreamUsecase(logDb))
		r.Get("/api/v1/streams/{stream}/fields", ListStreamFieldsUsecase(logDb))
		r.Delete("/api/v1/streams/{stream}", PurgeStreamUsecase(logDb))
	})

	return service
}

type routerWrapper struct {
	chi.Router
}

func (r *routerWrapper) Get(pattern string, u usecase.Interactor) {
	r.Method(http.MethodGet, pattern, nethttp.NewHandler(u))
}

func (r *routerWrapper) Post(pattern string, u usecase.Interactor) {
	r.Method(http.MethodPost, pattern, nethttp.NewHandler(u))
}

func (r *routerWrapper) Put(pattern string, u usecase.Interactor) {
	r.Method(http.MethodPut, pattern, nethttp.NewHandler(u))
}

func (r *routerWrapper) Delete(pattern string, u usecase.Interactor) {
	r.Method(http.MethodDelete, pattern, nethttp.NewHandler(u))
}
