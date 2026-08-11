package schemas

import (
	"link-society.com/flowg/internal/models"
)

// GetAuthProvidersRequest identifies the auth provider to retrieve.
type GetAuthProvidersRequest struct {
	// Name is the name of the auth provider to read.
	AuthProvider string `path:"auth_provider" minLength:"1"`
}

// GetAuthProvidersResponse carries the requested auth provider information.
type GetAuthProvidersResponse struct {
	// Success reports whether the auth provider was found and returned.
	Success bool `json:"success"`
	// AuthProvider is the auth provider and its details.
	AuthProvider *models.AuthProvider `json:"auth_provider"`
}
