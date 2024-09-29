package models

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"
	"net/http"
)

type WebhookV1 struct {
	Version int               `json:"version"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

func (w *WebhookV1) Call(ctx context.Context, record *LogRecord) error {
	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal log record: %w", err)
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
