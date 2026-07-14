package schemas

import (
	"link-society.com/flowg/internal/models"
)

// GetUserRequest identifies the user to retrieve.
type GetUserRequest struct {
	// Username is the name of the user to read.
	Username string `path:"user" minLength:"1"`
}

// GetUserResponse carries the requested user account.
type GetUserResponse struct {
	// Success reports whether the user was found and returned.
	Success bool `json:"success"`
	// User is the account and its assigned roles.
	User *models.User `json:"user"`
}
