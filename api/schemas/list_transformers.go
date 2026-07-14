package schemas

// ListTransformersRequest is empty: listing transformers takes no parameters.
type ListTransformersRequest struct{}

// ListTransformersResponse carries the names of the available transformers.
type ListTransformersResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Transformers holds the name of every configured transformer.
	Transformers []string `json:"transformers"`
}
