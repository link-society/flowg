package schemas

// ListForwardersRequest is empty: listing forwarders takes no parameters.
type ListForwardersRequest struct{}

// ListForwardersResponse carries the names of the available forwarders.
type ListForwardersResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Forwarders holds the name of every configured forwarder.
	Forwarders []string `json:"forwarders"`
}
