//go:build !debug

package mgmt

import "net/http"

// registerProfiler is a no-op in non-debug builds: the pprof endpoints are only
// mounted when the "debug" build tag is set (see profiler_debug.go).
func registerProfiler(*http.ServeMux) {
	// No-op in non-debug builds
}
