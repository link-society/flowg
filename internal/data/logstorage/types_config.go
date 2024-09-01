package logstorage

type StreamConfig struct {
	RetentionTime int64    `json:"ttl" description:"TTL in seconds, 0 to disable"`
	RetentionSize int64    `json:"size" description:"Maximum size in MB, 0 to disable"`
	IndexedFields []string `json:"indexed_fields"`
}

func (s StreamConfig) IsFieldIndexed(field string) bool {
	for _, f := range s.IndexedFields {
		if f == field {
			return true
		}
	}
	return false
}
