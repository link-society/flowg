package web

import (
	"net/http"

	"go.uber.org/fx"
)

// Options configures the web UI handler.
type Options struct {
	// MountPath is the path prefix the UI is served under (without trailing
	// slash), injected into the page so the UI can build its URLs.
	MountPath string
}

// Module wires the web UI handler into the application's fx dependency graph.
//
// It provides the http.Handler returned by NewHandler under the given
// handlerName so the HTTP server can mount it, configured by the [Options]
// provided to the graph.
func Module(handlerName string) fx.Option {
	return fx.Module(
		"web",
		fx.Provide(fx.Annotate(
			func(opts Options) http.Handler {
				return NewHandler(opts.MountPath)
			},
			fx.ResultTags(`name:"`+handlerName+`"`),
		)),
	)
}
