package alerting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"link-society.com/flowg/internal/data/logstorage"
)

type Webhook struct {
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

func (w *Webhook) Call(ctx context.Context, logEntry *logstorage.LogEntry) error {
	payload, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	buffer := bytes.NewBuffer(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", w.Url, buffer)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range w.Headers {
		req.Header.Add(key, value)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
