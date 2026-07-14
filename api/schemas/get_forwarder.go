package schemas

import "link-society.com/flowg/internal/models"

// GetForwarderRequest identifies the forwarder to retrieve.
type GetForwarderRequest struct {
	// Forwarder is the name of the forwarder to read.
	Forwarder string `path:"forwarder" minLength:"1"`
}

// GetForwarderResponse carries the definition of the requested forwarder.
type GetForwarderResponse struct {
	// Success reports whether the forwarder was found and returned.
	Success bool `json:"success"`
	// Forwarder is the forwarder's configuration.
	Forwarder *models.ForwarderV2 `json:"forwarder"`
}
