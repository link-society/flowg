package featureflags

import (
	"sync/atomic"
)

var (
	demoModeValue atomic.Bool
)

func GetDemoMode() bool {
	return demoModeValue.Load()
}

func SetDemoMode(enabled bool) {
	demoModeValue.Store(enabled)
}
