package schemas

// SaveRoleRequest carries the role name and the permissions to grant it.
type SaveRoleRequest struct {
	// Role is the name of the role to create or overwrite.
	Role string `path:"role" minLength:"1"`
	// Scopes are the names of the permissions to grant the role.
	Scopes []string `json:"scopes" required:"true"`
}

// SaveRoleResponse reports the outcome of the save.
type SaveRoleResponse struct {
	// Success reports whether the role was persisted.
	Success bool `json:"success"`
}
