package schemas

// PurgeStreamRequest identifies the stream to purge.
type PurgeStreamRequest struct {
	// Stream is the name of the stream to purge.
	Stream string `path:"stream" minLength:"1"`
}

// PurgeStreamResponse reports the outcome of the purge.
type PurgeStreamResponse struct {
	// Success reports whether the stream was purged.
	Success bool `json:"success"`
}
