package api

import (
	"go.uber.org/fx"

	_ "link-society.com/flowg/api/middlewares"
	_ "link-society.com/flowg/api/operations"
	"link-society.com/flowg/api/routing"
)

// Module bundles everything needed to serve FlowG's REST API into a single
// dependency-injection module.
//
// It pulls in the route table along with the operations and middlewares that
// populate it, so a consumer enables the whole API by depending on this one
// module instead of importing each endpoint package by hand. The handler is
// provided under the given name so callers can tell it apart from the other
// HTTP handlers they mount.
func Module(handlerName string) fx.Option {
	return fx.Module(
		"api",
		routing.Module(),
		fx.Provide(fx.Annotate(
			NewHandler,
			fx.ResultTags(`name:"`+handlerName+`"`),
		)),
	)
}
