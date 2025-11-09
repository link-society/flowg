package models

import (
	"context"
	"fmt"

	"bytes"
	"io"
	"net/http"

	otlpcommon "go.opentelemetry.io/proto/otlp/common/v1"
	otlplogs "go.opentelemetry.io/proto/otlp/logs/v1"
	proto "google.golang.org/protobuf/proto"
)

type forwarderStateOtlpV2 struct {
	client *http.Client
}

type ForwarderOtlpV2 struct {
	Type     string            `json:"type" enum:"otlp" required:"true"`
	Endpoint string            `json:"endpoint" required:"true"`
	Headers  map[string]string `json:"headers,omitempty"`

	state *forwarderStateOtlpV2
}

func (f *ForwarderOtlpV2) init(ctx context.Context) error {
	f.state = &forwarderStateOtlpV2{
		client: &http.Client{},
	}
	return nil
}

func (f *ForwarderOtlpV2) close(context.Context) error {
	return nil
}

func (f *ForwarderOtlpV2) call(ctx context.Context, record *LogRecord) error {
	logData := &otlplogs.LogsData{
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

	body, ok := record.Fields["body"]
	if !ok {
		body = ""
	}

	lr := &otlplogs.LogRecord{
		TimeUnixNano: uint64(record.Timestamp.UnixNano()),
		Body: &otlpcommon.AnyValue{
			Value: &otlpcommon.AnyValue_StringValue{
				StringValue: body,
			},
		},
		Attributes: []*otlpcommon.KeyValue{},
	}

	for k, v := range record.Fields {
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

	logData.ResourceLogs[0].ScopeLogs[0].LogRecords = append(logData.ResourceLogs[0].ScopeLogs[0].LogRecords, lr)

	data, err := proto.Marshal(logData)
	if err != nil {
		return fmt.Errorf("failed to marshal OTLP logs: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", f.Endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-protobuf")
	for k, v := range f.Headers {
		req.Header.Set(k, v)
	}

	resp, err := f.state.client.Do(req)
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
