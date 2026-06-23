package routing

import (
	"go.uber.org/fx"

	"github.com/swaggest/openapi-go"
	"github.com/swaggest/usecase"
)

// OperationsGroup is the dependency-injection value group every [Operation] is
// contributed to and that the HTTP handler consumes.
//
// It is documented here as the single source of truth, but the matching struct
// and result tags must spell it out as string literals, which Go requires.
const OperationsGroup = "operations"

// Operation describes a single API endpoint: the route that reaches it, the
// interactor that implements it, and how it is exposed.
//
// Keeping the route metadata next to the behaviour lets an endpoint be added,
// moved or removed by editing a single file.
type Operation struct {
	// Method is the HTTP method the route responds to.
	Method string
	// Pattern is the route pattern, including any path parameters.
	Pattern string
	// Public marks the operation as reachable without authentication; all other
	// operations sit behind the authentication middleware.
	Public bool
	// Interactor carries the operation's behaviour.
	Interactor usecase.Interactor
	// Annotate, when set, refines the operation's generated OpenAPI
	// documentation for cases the reflector cannot infer on its own.
	Annotate func(oc openapi.OperationContext) error
}

// OperationOption tweaks an [Operation] at registration time, expressing the
// traits that only some endpoints have (being public, carrying an OpenAPI
// annotation) without burdening the common case.
type OperationOption func(*Operation)

// Public marks the registered operation as reachable without authentication.
func Public() OperationOption {
	return func(op *Operation) { op.Public = true }
}

// Annotated attaches an OpenAPI annotation to the registered operation, for the
// few endpoints whose documentation the reflector cannot infer on its own.
func Annotated(annotate func(oc openapi.OperationContext) error) OperationOption {
	return func(op *Operation) { op.Annotate = annotate }
}

// RegisterOperation records an endpoint so [Module] provides it to the
// dependency-injection container's operations group.
//
// It binds a constructor to its route, injecting the constructor's own
// dependency struct, so the route metadata and the behaviour stay together.
// Operation files call it from an init function.
func RegisterOperation[Deps any](
	construct func(Deps) usecase.Interactor,
	method string,
	pattern string,
	opts ...OperationOption,
) {
	providers = append(providers, fx.Provide(fx.Annotate(
		func(deps Deps) Operation {
			op := Operation{Method: method, Pattern: pattern}
			for _, opt := range opts {
				opt(&op)
			}
			op.Interactor = construct(deps)
			return op
		},
		fx.ResultTags(`group:"operations"`),
	)))
}
