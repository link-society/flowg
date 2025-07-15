package models

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	otlpcommon "go.opentelemetry.io/proto/otlp/common/v1"
	otlplogs "go.opentelemetry.io/proto/otlp/logs/v1"
	proto "google.golang.org/protobuf/proto"
)

type OtlpForwarderConfig struct {
	Endpoint string            `json:"endpoint,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
}

type ForwarderOtlpV2 struct {
	Type   string              `json:"type" enum:"otlp"`
	Config OtlpForwarderConfig `json:"config"`
	client *http.Client        `json:"-"`
}

// call sends a single log record
func (f *ForwarderOtlpV2) call(ctx context.Context, record *LogRecord) error {
	// Initialize client if nil
	if f.client == nil {
		f.client = &http.Client{}
	}

	// Convert single log to OTLP protobuf msg format
	otlpLogs, err := ConvertToOtlpLogs([]*LogRecord{record})
	if err != nil {
		return err
	}

	// Marshal the OTLP logs to protobuf bytes
	data, err := proto.Marshal(otlpLogs)
	if err != nil {
		return fmt.Errorf("failed to marshal OTLP logs: %w", err)
	}

	// Create HTTP POST request
	req, err := http.NewRequestWithContext(ctx, "POST", f.Config.Endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set Headers
	req.Header.Set("Content-Type", "application/x-protobuf")
	for k, v := range f.Config.Headers {
		req.Header.Set(k, v)
	}

	// Send the request
	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ConvertToOtlpLogs converts a slice of LogRecord to OTLP LogsData format
func ConvertToOtlpLogs(records []*LogRecord) (*otlplogs.LogsData, error) {
	otlpLogs := &otlplogs.LogsData{
		ResourceLogs: []*otlplogs.ResourceLogs{
			{
				ScopeLogs: []*otlplogs.ScopeLogs{
					{
						LogRecords: []*otlplogs.LogRecord{},
					},
				},
			},
		},
	}

	for _, r := range records {
		lr := &otlplogs.LogRecord{
			TimeUnixNano: uint64(r.Timestamp.UnixNano()),
			Body: &otlpcommon.AnyValue{
				Value: &otlpcommon.AnyValue_StringValue{
					StringValue: getBody(r.Fields),
				},
			},
			Attributes: []*otlpcommon.KeyValue{},
		}

		// Convert log fields to OTLP attributes
		for k, v := range r.Fields {
			if k == "body" {
				continue // Skip body as it's already set
			}
			attr := &otlpcommon.KeyValue{
				Key: k,
				Value: &otlpcommon.AnyValue{
					Value: &otlpcommon.AnyValue_StringValue{
						StringValue: v,
					},
				},
			}
			lr.Attributes = append(lr.Attributes, attr)
		}

		otlpLogs.ResourceLogs[0].ScopeLogs[0].LogRecords = append(otlpLogs.ResourceLogs[0].ScopeLogs[0].LogRecords, lr)
	}

	return otlpLogs, nil
}

// getBody returns the body from fields or a default empty string
func getBody(fields map[string]string) string {
	body, ok := fields["body"]
	if !ok || body == "" {
		return ""
	}
	return body
}
