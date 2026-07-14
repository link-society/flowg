package schemas

import (
	"time"

	"link-society.com/flowg/internal/models"
)

// QueryStreamRequest describes a bounded search over a stream's logs.
type QueryStreamRequest struct {
	// Stream is the name of the stream to query.
	Stream string `path:"stream" minLength:"1"`
	// From is the inclusive lower bound of the time range.
	From time.Time `query:"from" format:"date-time" required:"true"`
	// To is the inclusive upper bound of the time range.
	To time.Time `query:"to" format:"date-time" required:"true"`
	// Filter is an optional filtering expression to match records against.
	Filter *string `query:"filter"`
	// Indexing narrows the search to specific values of indexed fields.
	Indexing map[string][]string `query:"indexing" collectionFormat:"json"`
}

// QueryStreamResponse carries the records matching the query.
type QueryStreamResponse struct {
	// Success reports whether the query completed.
	Success bool `json:"success"`
	// Records holds the matching log records.
	Records []models.LogRecord `json:"records"`
}
