package schemas

import (
	"link-society.com/flowg/internal/models"
)

// GetStreamRequest identifies the stream whose configuration is requested.
type GetStreamRequest struct {
	// Stream is the name of the stream to read.
	Stream string `path:"stream" minLength:"1"`
}

// GetStreamResponse carries the configuration of the requested stream.
type GetStreamResponse struct {
	// Success reports whether the configuration was returned.
	Success bool `json:"success"`
	// Config is the stream's retention and indexing configuration.
	Config models.StreamConfig `json:"config"`
}
