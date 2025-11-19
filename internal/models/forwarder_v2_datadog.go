package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type forwarderStateDatadogV2 struct {
	client *http.Client
}

type ForwarderDatadogV2 struct {
	Type   string `json:"type" enum:"datadog" required:"true"`
	Url    string `json:"url" required:"true" format:"uri"`
	ApiKey string `json:"apiKey" required:"true" minLength:"1"`

	state *forwarderStateDatadogV2
}

func (f *ForwarderDatadogV2) init(ctx context.Context) error {
	f.state = &forwarderStateDatadogV2{
		client: &http.Client{},
	}
	return nil
}

func (f *ForwarderDatadogV2) close(context.Context) error {
	return nil
}

func (f *ForwarderDatadogV2) call(ctx context.Context, record *LogRecord) error {
	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal log record: %w", err)
	}

	buffer := bytes.NewBuffer(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", f.Url, buffer)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("DD-API-KEY", f.ApiKey)

	resp, err := f.state.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
