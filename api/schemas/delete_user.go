package schemas

// DeleteUserRequest identifies the user to remove.
type DeleteUserRequest struct {
	// User is the name of the account to delete.
	User string `path:"user" minLength:"1"`
}

// DeleteUserResponse reports the outcome of the deletion.
type DeleteUserResponse struct {
	// Success reports whether the user was removed.
	Success bool `json:"success"`
}
