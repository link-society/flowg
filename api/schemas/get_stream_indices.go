package schemas

// GetStreamIndicesRequest identifies the stream whose index values are
// requested.
type GetStreamIndicesRequest struct {
	// Stream is the name of the stream to inspect.
	Stream string `path:"stream" minLength:"1"`
}

// GetStreamIndicesResponse carries the distinct values of each indexed field.
type GetStreamIndicesResponse struct {
	// Success reports whether the index values were returned.
	Success bool `json:"success"`
	// Indices maps each indexed field name to its known distinct values.
	Indices map[string][]string `json:"indices"`
}
