package schemas

import "link-society.com/flowg/internal/models"

// SaveForwarderRequest carries the forwarder name and its new definition.
type SaveForwarderRequest struct {
	// Forwarder is the name of the forwarder to create or overwrite.
	Forwarder string `path:"forwarder" minLength:"1"`
	// Config is the forwarder definition to store under that name.
	Config models.ForwarderV2 `json:"forwarder" required:"true"`
}

// SaveForwarderResponse reports the outcome of the save.
type SaveForwarderResponse struct {
	// Success reports whether the forwarder was persisted.
	Success bool `json:"success"`
}
