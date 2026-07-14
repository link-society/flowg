package schemas

import "link-society.com/flowg/internal/models"

// ListUsersRequest is empty: listing users takes no parameters.
type ListUsersRequest struct{}

// ListUsersResponse carries every known user with its roles.
type ListUsersResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Users holds every account and its assigned roles.
	Users []models.User `json:"users"`
}
