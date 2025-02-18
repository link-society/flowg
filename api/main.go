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

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

func NewHandler(
	authStorage *auth.Storage,
	logStorage *log.Storage,
	configStorage *config.Storage,
	logNotifier *lognotify.LogNotifier,
	pipelineRunner *pipelines.Runner,
) http.Handler {
	reflector := openapi31.NewReflector()
	service := web.NewService(reflector)

	service.OpenAPISchema().SetTitle("Flowg API")
	service.OpenAPISchema().SetVersion(app.FLOWG_VERSION)
	service.Docs("/api/docs", v5emb.New)

	service.Post("/api/v1/auth/login", LoginUsecase(authStorage))

	service.With(
		nethttp.HTTPBearerSecurityMiddleware(
			service.OpenAPICollector,
			"patAuth",
			"Authentication using Personal Access Token",
			"PAT",
		),
		nethttp.HTTPBearerSecurityMiddleware(
			service.OpenAPICollector,
			"jwtAuth",
			"Authentication using JSON Web Token",
			"JWT",
		),
		apiUtils.ApiMiddleware(authStorage),
	).Group(func(router chi.Router) {
		r := &routerWrapper{Router: router}

		r.Get("/api/v1/transformers", ListTransformersUsecase(authStorage, configStorage))
		r.Get("/api/v1/transformers/{transformer}", GetTransformerUsecase(authStorage, configStorage))
		r.Put("/api/v1/transformers/{transformer}", SaveTransformerUsecase(authStorage, configStorage))
		r.Delete("/api/v1/transformers/{transformer}", DeleteTransformerUsecase(authStorage, configStorage))

		r.Get("/api/v1/pipelines", ListPipelinesUsecase(authStorage, configStorage))
		r.Get("/api/v1/pipelines/{pipeline}", GetPipelineUsecase(authStorage, configStorage))
		r.Put("/api/v1/pipelines/{pipeline}", SavePipelineUsecase(authStorage, configStorage))
		r.Delete("/api/v1/pipelines/{pipeline}", DeletePipelineUsecase(authStorage, configStorage))
		r.Post("/api/v1/pipelines/{pipeline}/logs", IngestLogUsecase(authStorage, pipelineRunner))

		r.Get("/api/v1/streams", ListStreamsUsecase(authStorage, logStorage))
		r.Get("/api/v1/streams/{stream}", GetStreamUsecase(authStorage, logStorage))
		r.Put("/api/v1/streams/{stream}", ConfigureStreamUsecase(authStorage, logStorage))
		r.Get("/api/v1/streams/{stream}/logs", QueryStreamUsecase(authStorage, logStorage))
		r.Get("/api/v1/streams/{stream}/fields", ListStreamFieldsUsecase(authStorage, logStorage))
		r.Delete("/api/v1/streams/{stream}", PurgeStreamUsecase(authStorage, logStorage))

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
		r.Get("/api/v1/streams/{stream}/logs/watch", WatchLogsUsecase(authStorage, logNotifier))

		r.Get("/api/v1/alerts", ListAlertsUsecase(authStorage, configStorage))
		r.Get("/api/v1/alerts/{alert}", GetAlertUsecase(authStorage, configStorage))
		r.Put("/api/v1/alerts/{alert}", SaveAlertUsecase(authStorage, configStorage))
		r.Delete("/api/v1/alerts/{alert}", DeleteAlertUsecase(authStorage, configStorage))

		r.Get("/api/v1/roles", ListRolesUsecase(authStorage))
		r.Put("/api/v1/roles/{role}", SaveRoleUsecase(authStorage))
		r.Delete("/api/v1/roles/{role}", DeleteRoleUsecase(authStorage))

		r.Get("/api/v1/users", ListUsersUsecase(authStorage))
		r.Put("/api/v1/users/{user}", SaveUserUsecase(authStorage))
		r.Delete("/api/v1/users/{user}", DeleteUserUsecase(authStorage))

		r.Get("/api/v1/auth/whoami", WhoamiUsecase(authStorage))
		r.Post("/api/v1/auth/change-password", ChangePasswordUsecase(authStorage))

		r.Get("/api/v1/tokens", ListTokensUsecase(authStorage))
		r.Post("/api/v1/token", CreateTokenUsecase(authStorage))
		r.Delete("/api/v1/tokens/{token-uuid}", DeleteTokenUsecase(authStorage))

		r.Post("/api/v1/test/transformer", TestTransformerUsecase(authStorage))
		r.Post("/api/v1/test/alerts/{alert}", TestAlertUsecase(authStorage, configStorage))

		service.OpenAPICollector.AnnotateOperation(
			"GET", "/api/v1/backup/auth",
			func(oc openapi.OperationContext) error {
				contentUnits := oc.Response()
				for i, cu := range contentUnits {
					if cu.HTTPStatus == 200 {
						cu.ContentType = "application/octet-stream"
						cu.Description = "Binary file"
						cu.Format = "Binary file"
					}

					contentUnits[i] = cu
				}

				return nil
			},
		)
		r.Get("/api/v1/backup/auth", BackupAuthUsecase(authStorage))

		service.OpenAPICollector.AnnotateOperation(
			"GET", "/api/v1/backup/logs",
			func(oc openapi.OperationContext) error {
				contentUnits := oc.Response()
				for i, cu := range contentUnits {
					if cu.HTTPStatus == 200 {
						cu.ContentType = "application/octet-stream"
						cu.Description = "Binary file"
						cu.Format = "Binary file"
					}

					contentUnits[i] = cu
				}

				return nil
			},
		)
		r.Get("/api/v1/backup/logs", BackupLogsUsecase(authStorage, logStorage))

		service.OpenAPICollector.AnnotateOperation(
			"GET", "/api/v1/backup/config",
			func(oc openapi.OperationContext) error {
				contentUnits := oc.Response()
				for i, cu := range contentUnits {
					if cu.HTTPStatus == 200 {
						cu.ContentType = "application/octet-stream"
						cu.Description = "Binary file"
						cu.Format = "Binary file"
					}

					contentUnits[i] = cu
				}

				return nil
			},
		)
		r.Get("/api/v1/backup/config", BackupConfigUsecase(authStorage, configStorage))

		r.Post("/api/v1/restore/auth", RestoreAuthUsecase(authStorage))
		r.Post("/api/v1/restore/logs", RestoreLogsUsecase(authStorage, logStorage))
		r.Post("/api/v1/restore/config", RestoreConfigUsecase(authStorage, configStorage))
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
