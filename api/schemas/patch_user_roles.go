package schemas

// PatchUserRolesRequest carries the user and the roles to assign it.
type PatchUserRolesRequest struct {
	// User is the name of the account to update.
	User string `path:"user" minLength:"1"`
	// Roles are the names of the roles to assign, replacing the current set.
	Roles []string `json:"roles" required:"true" items.minLength:"1"`
}

// PatchUserRolesResponse reports the outcome of the update.
type PatchUserRolesResponse struct {
	// Success reports whether the roles were updated.
	Success bool `json:"success"`
}
