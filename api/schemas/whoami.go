package schemas

import (
	"link-society.com/flowg/internal/models"
)

// WhoamiRequest is empty: the caller is identified by their credentials.
type WhoamiRequest struct{}

// WhoamiResponse describes the currently authenticated user.
type WhoamiResponse struct {
	// Success reports whether the profile was returned.
	Success bool `json:"success"`
	// User is the authenticated account and its assigned roles.
	User *models.User `json:"user"`
	// Permissions is the effective permission set derived from the user's roles.
	Permissions models.Permissions `json:"permissions"`
}
