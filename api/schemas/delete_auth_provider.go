package schemas

// DeleteAuthProviderRequest identifies the auth provider to remove.
type DeleteAuthProviderRequest struct {
	// AuthProvider is the name of the auth provider to delete.
	AuthProvider string `path:"auth_provider" minLength:"1"`
}

// DeleteAuthProviderResponse reports the outcome of the deletion.
type DeleteAuthProviderResponse struct {
	// Success reports whether the auth provider was removed.
	Success bool `json:"success"`
}
