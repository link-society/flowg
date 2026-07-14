package schemas

import "link-society.com/flowg/internal/models"

// ListStreamsRequest is empty: listing streams takes no parameters.
type ListStreamsRequest struct{}

// ListStreamsResponse carries every known stream and its configuration.
type ListStreamsResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Streams maps each stream name to its configuration.
	Streams map[string]models.StreamConfig `json:"streams"`
}
