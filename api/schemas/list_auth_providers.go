package schemas

// AuthProviderInfo carries the name and display name of an auth provider.
type AuthProviderInfo struct {
	Name string `json:"name" required:"true"`
	DisplayName string `json:"display_name" required:"true"`
}

// ListAuthProvidersRequest is empty: listing auth providers takes no parameters.
type ListAuthProvidersRequest struct{}

// ListAuthProvidersResponse carries the names and display names of the available auth providers.
type ListAuthProvidersResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// AuthProviders holds the name and display name of every configured auth provider.
	AuthProviders []AuthProviderInfo `json:"auth_providers"`
}

