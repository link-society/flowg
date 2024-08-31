package logstorage

import "time"

type StreamConfig struct {
	RetentionTime time.Duration `json:"ttl"`
	RetentionSize int64         `json:"size"`
}
