package cmd

// ContextKey is the type of the keys under which the command context carries the
// configured API clients.
type ContextKey string

const (
	// ApiClient is the context key for the client targeting the FlowG HTTP API.
	ApiClient ContextKey = "api_client"
	// MgmtApiClient is the context key for the client targeting the FlowG
	// management HTTP API.
	MgmtApiClient ContextKey = "mgmt_api_client"
)
