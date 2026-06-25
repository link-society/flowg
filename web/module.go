package web

import (
	"net/http"

	"go.uber.org/fx"
)

// Module wires the web UI handler into the application's fx dependency graph.
//
// It provides the http.Handler returned by NewHandler under the given
// handlerName so the HTTP server can mount it, and configures the handler to
// serve the UI from mountPath.
func Module(handlerName string, mountPath string) fx.Option {
	return fx.Module(
		"web",
		fx.Provide(fx.Annotate(
			func() http.Handler {
				return NewHandler(mountPath)
			},
			fx.ResultTags(`name:"`+handlerName+`"`),
		)),
	)
}
