//go:build !debug

package mgmt

import "net/http"

func registerProfiler(*http.ServeMux) {
	// No-op in non-debug builds
}
