package api

import (
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	"github.com/swaggest/swgui/v5emb"

	"go.uber.org/fx"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/middlewares"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/app"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"

	authStorage "link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

// handlerParams gathers everything [NewHandler] needs from the
// dependency-injection container: the cross-cutting backends the handler wires
// directly, and the full set of operations contributed to the routing group.
type handlerParams struct {
	fx.In

	AuthStorage   authStorage.Storage
	LogStorage    log.Storage
	ConfigStorage config.Storage

	LogNotifier    lognotify.LogNotifier
	PipelineRunner pipelines.Runner

	Operations []routing.Operation `group:"operations"`
}

// NewHandler builds the HTTP handler that serves FlowG's REST API.
//
// It owns the concerns that surround the operations themselves: the OpenAPI
// service and its documentation, the routing table that maps each operation to
// its route, and the security middlewares that authenticate callers. The
// behaviour and route of every endpoint live in the [operations] package,
// keeping this file focused on assembly rather than enumeration.
func NewHandler(params handlerParams) http.Handler {
	reflector := openapi31.NewReflector()
	service := web.NewService(reflector)

	service.OpenAPISchema().SetTitle("Flowg API")
	service.OpenAPISchema().SetVersion(app.FLOWG_VERSION)
	service.Docs("/api/docs", v5emb.New)

	// Register operations in a stable order. The OpenAPI collector derives each
	// shared component schema from the first operation that references the
	// underlying type, so a deterministic order keeps the generated document
	// reproducible regardless of the order the container yields the group.
	ops := slices.Clone(params.Operations)
	slices.SortFunc(ops, func(a, b routing.Operation) int {
		if c := strings.Compare(a.Pattern, b.Pattern); c != 0 {
			return c
		}
		return strings.Compare(a.Method, b.Method)
	})

	register := func(router interface {
		Method(method, pattern string, h http.Handler)
	}, op routing.Operation) {
		if op.Annotate != nil {
			service.OpenAPICollector.AnnotateOperation(op.Method, op.Pattern, op.Annotate)
		}

		router.Method(op.Method, op.Pattern, nethttp.NewHandler(op.Interactor))
	}

	for _, op := range ops {
		if op.Public {
			register(service, op)
		}
	}

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
		auth.ApiMiddleware(params.AuthStorage),
	).Group(func(router chi.Router) {
		for _, op := range ops {
			if !op.Public {
				register(router, op)
			}
		}
	})

	service.Mount("/api/v1/middlewares/", middlewares.NewHandler(&middlewares.Dependencies{
		AuthStorage:   params.AuthStorage,
		LogStorage:    params.LogStorage,
		ConfigStorage: params.ConfigStorage,

		LogNotifier:    params.LogNotifier,
		PipelineRunner: params.PipelineRunner,
	}))

	return service
}
