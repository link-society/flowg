package sse

// Event is a single decoded Server-Sent Event.
type Event struct {
	ID   string
	Type string
	Data string
}
