package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	"github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"

	"link-society.com/flowg/internal/app"
	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/data/lognotify"
	"link-society.com/flowg/internal/data/logstorage"
)

func NewHandler(
	authDb *auth.Database,
	logDb *logstorage.Storage,
	configStorage *config.Storage,
	logNotifier *lognotify.LogNotifier,
) http.Handler {
	reflector := openapi31.NewReflector()
	service := web.NewService(reflector)

	service.OpenAPISchema().SetTitle("Flowg API")
	service.OpenAPISchema().SetVersion(app.FLOWG_VERSION)
	service.Docs("/api/docs", v5emb.New)

	service.Post("/api/v1/auth/login", LoginUsecase(authDb))

	service.With(
		nethttp.HTTPBasicSecurityMiddleware(
			service.OpenAPICollector,
			"patAuth",
			"Authentication using Personal Access Token",
		),
		nethttp.HTTPBearerSecurityMiddleware(
			service.OpenAPICollector,
			"jwtAuth",
			"Authentication using JSON Web Token",
			"JWT",
		),
		auth.ApiMiddleware(authDb),
	).Group(func(router chi.Router) {
		r := &routerWrapper{Router: router}

		r.Get("/api/v1/transformers", ListTransformersUsecase(authDb, configStorage))
		r.Get("/api/v1/transformers/{transformer}", GetTransformerUsecase(authDb, configStorage))
		r.Put("/api/v1/transformers/{transformer}", SaveTransformerUsecase(authDb, configStorage))
		r.Delete("/api/v1/transformers/{transformer}", DeleteTransformerUsecase(authDb, configStorage))

		r.Get("/api/v1/pipelines", ListPipelinesUsecase(authDb, configStorage))
		r.Get("/api/v1/pipelines/{pipeline}", GetPipelineUsecase(authDb, configStorage))
		r.Put("/api/v1/pipelines/{pipeline}", SavePipelineUsecase(authDb, configStorage))
		r.Delete("/api/v1/pipelines/{pipeline}", DeletePipelineUsecase(authDb, configStorage))
		r.Post("/api/v1/pipelines/{pipeline}/logs", IngestLogUsecase(authDb, configStorage, logDb, logNotifier))

		r.Get("/api/v1/streams", ListStreamsUsecase(authDb, logDb))
		r.Get("/api/v1/streams/{stream}", GetStreamUsecase(authDb, logDb))
		r.Put("/api/v1/streams/{stream}", ConfigureStreamUsecase(authDb, logDb))
		r.Get("/api/v1/streams/{stream}/logs", QueryStreamUsecase(authDb, logDb))
		r.Get("/api/v1/streams/{stream}/fields", ListStreamFieldsUsecase(authDb, logDb))
		r.Delete("/api/v1/streams/{stream}", PurgeStreamUsecase(authDb, logDb))

		service.OpenAPICollector.AnnotateOperation(
			"GET", "/api/v1/streams/{stream}/logs/watch",
			func(oc openapi.OperationContext) error {
				contentUnits := oc.Response()
				for i, cu := range contentUnits {
					if cu.HTTPStatus == 200 {
						cu.ContentType = "text/event-stream"
						cu.Description = "Stream of log entries"
						cu.Format = "Server-Sent Events"
					}

					contentUnits[i] = cu
				}

				return nil
			},
		)
		r.Get("/api/v1/streams/{stream}/logs/watch", WatchLogsUsecase(authDb, logNotifier))

		r.Get("/api/v1/alerts", ListAlertsUsecase(authDb, configStorage))
		r.Get("/api/v1/alerts/{alert}", GetAlertUsecase(authDb, configStorage))
		r.Put("/api/v1/alerts/{alert}", SaveAlertUsecase(authDb, configStorage))
		r.Delete("/api/v1/alerts/{alert}", DeleteAlertUsecase(authDb, configStorage))

		r.Get("/api/v1/roles", ListRolesUsecase(authDb))
		r.Put("/api/v1/roles/{role}", SaveRoleUsecase(authDb))
		r.Delete("/api/v1/roles/{role}", DeleteRoleUsecase(authDb))

		r.Get("/api/v1/users", ListUsersUsecase(authDb))
		r.Put("/api/v1/users/{user}", SaveUserUsecase(authDb))
		r.Delete("/api/v1/users/{user}", DeleteUserUsecase(authDb))

		r.Get("/api/v1/auth/whoami", WhoamiUsecase(authDb))
		r.Post("/api/v1/auth/change-password", ChangePasswordUsecase(authDb))

		r.Get("/api/v1/tokens", ListTokensUsecase(authDb))
		r.Post("/api/v1/token", CreateTokenUsecase(authDb))
		r.Delete("/api/v1/tokens/{token-uuid}", DeleteTokenUsecase(authDb))

		r.Post("/api/v1/test/transformer", TestTransformerUsecase(authDb))
		r.Post("/api/v1/test/alerts/{alert}", TestAlertUsecase(authDb, configStorage))
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
