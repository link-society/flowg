package schemas

import "link-society.com/flowg/internal/models"

// GetRoleRequest identifies the role to retrieve.
type GetRoleRequest struct {
	// Role is the name of the role to read.
	Role string `path:"role" minLength:"1"`
}

// GetRoleResponse carries the definition of the requested role.
type GetRoleResponse struct {
	// Success reports whether the role was found and returned.
	Success bool `json:"success"`
	// Role is the role and its granted permissions.
	Role *models.Role `json:"role"`
}
