package models

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"

	"net/http"
)

type ForwarderV2 struct {
	Version int               `json:"version"`
	Config  ForwarderConfigV2 `json:"config"`
}

type ForwarderConfigV2 struct {
	Http *ForwarderHttpV2 `json:"-"`
}

type ForwarderHttpV2 struct {
	Type    string            `json:"type" enum:"http"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

func (f *ForwarderV2) Call(ctx context.Context, record *LogRecord) error {
	switch {
	case f.Config.Http != nil:
		return f.Config.Http.call(ctx, record)

	default:
		return fmt.Errorf("unsupported forwarder type")
	}
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

func (cfg ForwarderConfigV2) MarshalJSON() ([]byte, error) {
	switch {
	case cfg.Http != nil:
		return json.Marshal(cfg.Http)

	default:
		return nil, fmt.Errorf("unsupported forwarder type")
	}
}

func (cfg ForwarderConfigV2) UnmarshalJSON(data []byte) error {
	var typeInfo struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &typeInfo); err != nil {
		return fmt.Errorf("failed to unmarshal forwarder type: %w", err)
	}

	switch typeInfo.Type {
	case "http":
		cfg.Http = &ForwarderHttpV2{}
		return json.Unmarshal(data, cfg.Http)

	default:
		return fmt.Errorf("unsupported forwarder type: %s", typeInfo.Type)
	}
}

func (ForwarderConfigV2) JSONSchemaOneOf() []any {
	return []any{
		ForwarderHttpV2{},
	}
}
