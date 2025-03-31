package api

import (
	"log/slog"

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

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
)

type Dependencies struct {
	AuthStorage   *auth.Storage
	LogStorage    *log.Storage
	ConfigStorage *config.Storage

	LogNotifier    *lognotify.LogNotifier
	PipelineRunner *pipelines.Runner
}

type controller struct {
	logger *slog.Logger
	deps   *Dependencies
}

func NewHandler(deps *Dependencies) http.Handler {
	ctrl := &controller{
		logger: slog.Default().With(slog.String("channel", "api")),
		deps:   deps,
	}

	reflector := openapi31.NewReflector()
	service := web.NewService(reflector)

	service.OpenAPISchema().SetTitle("Flowg API")
	service.OpenAPISchema().SetVersion(app.FLOWG_VERSION)
	service.Docs("/api/docs", v5emb.New)

	service.Post("/api/v1/auth/login", ctrl.LoginUsecase())

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
		apiUtils.ApiMiddleware(deps.AuthStorage),
	).Group(func(router chi.Router) {
		r := &routerWrapper{Router: router}

		r.Get("/api/v1/transformers", ctrl.ListTransformersUsecase())
		r.Get("/api/v1/transformers/{transformer}", ctrl.GetTransformerUsecase())
		r.Put("/api/v1/transformers/{transformer}", ctrl.SaveTransformerUsecase())
		r.Delete("/api/v1/transformers/{transformer}", ctrl.DeleteTransformerUsecase())

		r.Get("/api/v1/pipelines", ctrl.ListPipelinesUsecase())
		r.Get("/api/v1/pipelines/{pipeline}", ctrl.GetPipelineUsecase())
		r.Put("/api/v1/pipelines/{pipeline}", ctrl.SavePipelineUsecase())
		r.Delete("/api/v1/pipelines/{pipeline}", ctrl.DeletePipelineUsecase())

		r.Post("/api/v1/pipelines/{pipeline}/logs/struct", ctrl.IngestLogsStructUsecase())
		r.Post("/api/v1/pipelines/{pipeline}/logs/text", ctrl.IngestLogsTextUsecase())
		r.Post("/api/v1/pipelines/{pipeline}/logs/otlp", ctrl.IngestLogsOTLPUsecase())

		r.Get("/api/v1/streams", ctrl.ListStreamsUsecase())
		r.Get("/api/v1/streams/{stream}", ctrl.GetStreamUsecase())
		r.Put("/api/v1/streams/{stream}", ctrl.ConfigureStreamUsecase())
		r.Get("/api/v1/streams/{stream}/logs", ctrl.QueryStreamUsecase())
		r.Get("/api/v1/streams/{stream}/fields", ctrl.ListStreamFieldsUsecase())
		r.Delete("/api/v1/streams/{stream}", ctrl.PurgeStreamUsecase())

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
		r.Get("/api/v1/streams/{stream}/logs/watch", ctrl.WatchLogsUsecase())

		r.Get("/api/v1/forwarders", ctrl.ListForwardersUsecase())
		r.Get("/api/v1/forwarders/{forwarder}", ctrl.GetForwarderUsecase())
		r.Put("/api/v1/forwarders/{forwarder}", ctrl.SaveForwarderUsecase())
		r.Delete("/api/v1/forwarders/{forwarder}", ctrl.DeleteForwarderUsecase())

		r.Get("/api/v1/roles", ctrl.ListRolesUsecase())
		r.Get("/api/v1/roles/{role}", ctrl.GetRoleUsecase())
		r.Put("/api/v1/roles/{role}", ctrl.SaveRoleUsecase())
		r.Delete("/api/v1/roles/{role}", ctrl.DeleteRoleUsecase())

		r.Get("/api/v1/users", ctrl.ListUsersUsecase())
		r.Get("/api/v1/users/{user}", ctrl.GetUserUsecase())
		r.Put("/api/v1/users/{user}", ctrl.SaveUserUsecase())
		r.Patch("/api/v1/users/{user}", ctrl.PatchUserRolesUsecase())
		r.Delete("/api/v1/users/{user}", ctrl.DeleteUserUsecase())

		r.Get("/api/v1/auth/whoami", ctrl.WhoamiUsecase())
		r.Post("/api/v1/auth/change-password", ctrl.ChangePasswordUsecase())

		r.Get("/api/v1/tokens", ctrl.ListTokensUsecase())
		r.Post("/api/v1/token", ctrl.CreateTokenUsecase())
		r.Delete("/api/v1/tokens/{token-uuid}", ctrl.DeleteTokenUsecase())

		r.Post("/api/v1/test/transformer", ctrl.TestTransformerUsecase())
		r.Post("/api/v1/test/forwarders/{forwarder}", ctrl.TestForwarderUsecase())

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
		r.Get("/api/v1/backup/auth", ctrl.BackupAuthUsecase())

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
		r.Get("/api/v1/backup/logs", ctrl.BackupLogsUsecase())

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
		r.Get("/api/v1/backup/config", ctrl.BackupConfigUsecase())

		r.Post("/api/v1/restore/auth", ctrl.RestoreAuthUsecase())
		r.Post("/api/v1/restore/logs", ctrl.RestoreLogsUsecase())
		r.Post("/api/v1/restore/config", ctrl.RestoreConfigUsecase())
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

func (r *routerWrapper) Patch(pattern string, u usecase.Interactor) {
	r.Method(http.MethodPatch, pattern, nethttp.NewHandler(u))
}

func (r *routerWrapper) Delete(pattern string, u usecase.Interactor) {
	r.Method(http.MethodDelete, pattern, nethttp.NewHandler(u))
}
