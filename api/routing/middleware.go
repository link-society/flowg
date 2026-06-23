package routing

import (
	"net/http"

	"go.uber.org/fx"
)

// MiddlewaresGroup is the dependency-injection value group every [Middleware] is
// contributed to and that the HTTP handler consumes.
//
// It is documented here as the single source of truth, but the matching struct
// and result tags must spell it out as string literals, which Go requires.
const MiddlewaresGroup = "middlewares"

// Middleware describes a protocol-compatibility sub-handler: the route prefix it
// is mounted under and the handler that serves it.
//
// Keeping the mount point next to the handler lets a middleware be added, moved
// or removed by editing a single file, the same way [Operation] does for
// endpoints.
type Middleware struct {
	// Pattern is the route prefix the middleware is mounted under.
	Pattern string
	// Handler serves every request below [Middleware.Pattern].
	Handler http.Handler
}

// RegisterMiddleware records a middleware so [Module] provides it to the
// dependency-injection container's middlewares group.
//
// It binds a constructor to its mount point, injecting the constructor's own
// dependency struct, so the route prefix and the handler stay together.
// Middleware files call it from an init function.
func RegisterMiddleware[Deps any](
	construct func(Deps) http.Handler,
	pattern string,
) {
	providers = append(providers, fx.Provide(fx.Annotate(
		func(deps Deps) Middleware {
			return Middleware{Pattern: pattern, Handler: construct(deps)}
		},
		fx.ResultTags(`group:"middlewares"`),
	)))
}
