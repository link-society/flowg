package routing

import (
	"go.uber.org/fx"
)

// providers accumulates one fx provider per registered endpoint, so [Module]
// can expose the whole set without a hand-maintained list.
var providers []fx.Option

// Module provides every registered operation to the dependency-injection
// container.
//
// It must be called after the operation files have registered themselves, which
// the Go runtime guarantees by running their init functions before any caller.
func Module() fx.Option {
	return fx.Module("api.routing", providers...)
}
