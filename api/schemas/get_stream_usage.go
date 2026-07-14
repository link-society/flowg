package schemas

// GetStreamUsageRequest identifies the stream whose disk usage is requested.
type GetStreamUsageRequest struct {
	// Stream is the name of the stream to measure.
	Stream string `path:"stream" minLength:"1"`
}

// GetStreamUsageResponse carries the measured storage footprint.
type GetStreamUsageResponse struct {
	// Success reports whether the measurement was returned.
	Success bool `json:"success"`
	// Usage is the storage footprint of the stream, in bytes.
	Usage int64 `json:"usage"`
}
