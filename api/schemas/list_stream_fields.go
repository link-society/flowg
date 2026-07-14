package schemas

// ListStreamFieldsRequest identifies the stream whose fields are requested.
type ListStreamFieldsRequest struct {
	// Stream is the name of the stream to inspect.
	Stream string `path:"stream" minLength:"1"`
}

// ListStreamFieldsResponse carries the field names observed in the stream.
type ListStreamFieldsResponse struct {
	// Success reports whether the field names were returned.
	Success bool `json:"success"`
	// Fields holds every field name seen across the stream's records.
	Fields []string `json:"fields"`
}
