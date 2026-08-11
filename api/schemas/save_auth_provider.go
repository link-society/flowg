package schemas

import "link-society.com/flowg/internal/models"

// SaveAuthProviderRequest carries the auth provider name and its new definition.
type SaveAuthProviderRequest struct {
	// AuthProvider is the name of the auth provider to create or overwrite.
	AuthProvider string `path:"auth_provider" minLength:"1"`
	// Config is the auth provider definition to store under that name.
	Config models.AuthProvider `json:"auth_provider" required:"true"`
}

// SaveAuthProviderResponse reports the outcome of the save.
type SaveAuthProviderResponse struct {
	// Success reports whether the auth provider was persisted.
	Success bool `json:"success"`
}
