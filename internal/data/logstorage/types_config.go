package logstorage

type StreamConfig struct {
	RetentionTime int64 `json:"ttl" description:"TTL in seconds"`
	RetentionSize int64 `json:"size" description:"Maximum size in MB"`
}
