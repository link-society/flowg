package schemas

// DeleteRoleRequest identifies the role to remove.
type DeleteRoleRequest struct {
	// Role is the name of the role to delete.
	Role string `path:"role" minLength:"1"`
}

// DeleteRoleResponse reports the outcome of the deletion.
type DeleteRoleResponse struct {
	// Success reports whether the role was removed.
	Success bool `json:"success"`
}
