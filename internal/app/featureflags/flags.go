package featureflags

// Flags holds the feature toggles enabled for this instance. It is provided to
// the dependency graph from the configuration.
type Flags struct {
	// DemoMode configures the instance for public demonstrations: it seeds a
	// "demo" role and user, enforces a short retention on every stream, and is
	// exposed to the web UI. It is not suitable for production use.
	DemoMode bool
}
