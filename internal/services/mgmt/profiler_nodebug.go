//go:build !debug

package mgmt

import "net/http"

func registerProfiler(*http.ServeMux) {}
