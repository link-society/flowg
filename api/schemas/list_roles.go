package schemas

import "link-society.com/flowg/internal/models"

// ListRolesRequest is empty: listing roles takes no parameters.
type ListRolesRequest struct{}

// ListRolesResponse carries every known role with its permissions.
type ListRolesResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Roles holds every configured role and its granted permissions.
	Roles []models.Role `json:"roles"`
}
