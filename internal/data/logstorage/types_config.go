package logstorage

type StreamConfig struct {
	RetentionTime int64 `json:"ttl"`
	RetentionSize int64 `json:"size"`
}
