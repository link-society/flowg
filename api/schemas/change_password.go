package schemas

// ChangePasswordRequest carries the caller's current and desired passwords.
type ChangePasswordRequest struct {
	// OldPassword is the caller's current password, required to authorize the
	// change.
	OldPassword string `json:"old_password" required:"true" minLength:"1"`
	// NewPassword is the password to set.
	NewPassword string `json:"new_password" required:"true" minLength:"1"`
}

// ChangePasswordResponse reports the outcome of the change.
type ChangePasswordResponse struct {
	// Success reports whether the password was updated.
	Success bool `json:"success"`
}
