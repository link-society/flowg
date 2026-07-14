package schemas

import "link-society.com/flowg/internal/models"

// ConfigureStreamRequest carries the stream name and its new configuration.
type ConfigureStreamRequest struct {
	// Stream is the name of the stream to configure.
	Stream string `path:"stream" minLength:"1"`
	// Config is the retention and indexing configuration to apply.
	Config models.StreamConfig `json:"config" required:"true"`
}

// ConfigureStreamResponse reports the outcome of the configuration change.
type ConfigureStreamResponse struct {
	// Success reports whether the configuration was applied.
	Success bool `json:"success"`
}
