package models

// StreamConfig is the per-stream retention and indexing policy. Retention can be
// bounded by age (RetentionTime) and/or size (RetentionSize), and IndexedFields
// lists the fields kept in the inverted index for fast lookups.
type StreamConfig struct {
	RetentionTime int64    `json:"ttl" required:"true" description:"TTL in seconds, 0 to disable"`
	RetentionSize int64    `json:"size" required:"true" description:"Maximum size in MB, 0 to disable"`
	IndexedFields []string `json:"indexed_fields" items.minLength:"1"`
}

// IsFieldIndexed reports whether the named field is part of the stream's
// inverted index.
func (s StreamConfig) IsFieldIndexed(field string) bool {
	for _, f := range s.IndexedFields {
		if f == field {
			return true
		}
	}
	return false
}
