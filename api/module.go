package api

import (
	"log/slog"

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
//
// Within the module, the root logger is decorated with the "api" channel, so
// every operation and middleware logs under it without further setup.
func Module(handlerName string) fx.Option {
	return fx.Module(
		"api",
		fx.Decorate(func(rootLogger *slog.Logger) *slog.Logger {
			return rootLogger.With(slog.String("channel", "api"))
		}),
		routing.Module(),
		fx.Provide(fx.Annotate(
			NewHandler,
			fx.ResultTags(`name:"`+handlerName+`"`),
		)),
	)
}
