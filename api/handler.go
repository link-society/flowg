package api

import (
	"net/http"
	"slices"
	"strings"

	"go.uber.org/fx"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	"github.com/swaggest/swgui/v5emb"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/app"

	authStorage "link-society.com/flowg/internal/storage/auth"
)

// handlerParams gathers everything [NewHandler] needs from the
// dependency-injection container: the authentication backend it guards the
// routes with, the protocol-compatibility middlewares it mounts, and the full
// set of operations contributed to the routing group.
type handlerParams struct {
	fx.In

	AuthStorage authStorage.Storage

	Operations  []routing.Operation  `group:"operations"`
	Middlewares []routing.Middleware `group:"middlewares"`
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

	for _, mw := range params.Middlewares {
		service.Mount(mw.Pattern, mw.Handler)
	}

	return service
}
