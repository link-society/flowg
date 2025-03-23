package models

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"

	"net/http"
)

type ForwarderHttpV2 struct {
	Type    string            `json:"type" enum:"http"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

func (f *ForwarderHttpV2) call(ctx context.Context, record *LogRecord) error {
	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal log record: %w", err)
	}

	buffer := bytes.NewBuffer(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", f.Url, buffer)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range f.Headers {
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
