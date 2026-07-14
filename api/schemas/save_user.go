package schemas

// SaveUserRequest carries a user account and its initial password.
type SaveUserRequest struct {
	// User is the name of the account to create or overwrite.
	User string `path:"user" minLength:"1"`
	// Roles are the names of the roles to assign the user.
	Roles []string `json:"roles" required:"true"`
	// Password is the account's password; it is stored hashed.
	Password string `json:"password" required:"true"`
}

// SaveUserResponse reports the outcome of the save.
type SaveUserResponse struct {
	// Success reports whether the user was persisted.
	Success bool `json:"success"`
}
