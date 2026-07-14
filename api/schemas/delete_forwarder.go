package schemas

// DeleteForwarderRequest identifies the forwarder to remove.
type DeleteForwarderRequest struct {
	// Forwarder is the name of the forwarder to delete.
	Forwarder string `path:"forwarder" minLength:"1"`
}

// DeleteForwarderResponse reports the outcome of the deletion.
type DeleteForwarderResponse struct {
	// Success reports whether the forwarder was removed.
	Success bool `json:"success"`
}
