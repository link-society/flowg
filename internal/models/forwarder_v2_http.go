package models

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"

	"net/http"
)

type forwarderStateHttpV2 struct {
	client *http.Client
}

type ForwarderHttpV2 struct {
	Type    string            `json:"type" enum:"http" required:"true"`
	Url     string            `json:"url" required:"true"`
	Headers map[string]string `json:"headers,omitempty"`

	state *forwarderStateHttpV2
}

func (f *ForwarderHttpV2) init(context.Context) error {
	f.state = &forwarderStateHttpV2{
		client: &http.Client{},
	}
	return nil
}

func (f *ForwarderHttpV2) close(context.Context) error {
	return nil
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

	resp, err := f.state.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
