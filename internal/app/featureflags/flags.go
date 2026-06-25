package featureflags

import (
	"sync/atomic"
)

var (
	demoModeValue atomic.Bool
)

// GetDemoMode reports whether demo mode is enabled. In demo mode FlowG disables
// state-mutating operations so a public instance can be exposed safely.
func GetDemoMode() bool {
	return demoModeValue.Load()
}

// SetDemoMode enables or disables demo mode.
func SetDemoMode(enabled bool) {
	demoModeValue.Store(enabled)
}
