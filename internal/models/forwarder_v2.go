package models

import (
	"context"
	"fmt"

	"encoding/json"
)

// ForwarderV2 is the current forwarder model: a destination a router/forward
// node can send records to. The concrete destination lives in Config, which is a
// tagged union of one backend type.
type ForwarderV2 struct {
	Version int               `json:"version" default:"2"`
	Config  ForwarderConfigV2 `json:"config" required:"true"`
}

// ForwarderConfigV2 is a tagged union: exactly one field is non-nil, selecting
// the forwarder backend. It marshals to/from the backend's own JSON (discriminated
// by a "type" field) rather than nesting under the field name.
type ForwarderConfigV2 struct {
	Http               *ForwarderHttpV2               `json:"-"`
	Syslog             *ForwarderSyslogV2             `json:"-"`
	Datadog            *ForwarderDatadogV2            `json:"-"`
	Amqp               *ForwarderAmqpV2               `json:"-"`
	Splunk             *ForwarderSplunkV2             `json:"-"`
	Otlp               *ForwarderOtlpV2               `json:"-"`
	Elastic            *ForwarderElasticV2            `json:"-"`
	Clickhouse         *ForwarderClickhouseV2         `json:"-"`
	CloudWatch         *ForwarderCloudWatchV2         `json:"-"`
	GoogleCloudLogging *ForwarderGoogleCloudLoggingV2 `json:"-"`
}

// JSONSchemaOneOf advertises every backend variant so the generated OpenAPI
// schema models Config as a "oneOf".
func (ForwarderConfigV2) JSONSchemaOneOf() []any {
	return []any{
		ForwarderHttpV2{},
		ForwarderSyslogV2{},
		ForwarderDatadogV2{},
		ForwarderAmqpV2{},
		ForwarderSplunkV2{},
		ForwarderOtlpV2{},
		ForwarderElasticV2{},
		ForwarderClickhouseV2{},
		ForwarderCloudWatchV2{},
		ForwarderGoogleCloudLoggingV2{},
	}
}

// Init prepares the active backend (e.g. opens a connection pool) before any
// record is forwarded.
func (f *ForwarderV2) Init(ctx context.Context) error {
	switch {
	case f.Config.Http != nil:
		return f.Config.Http.init(ctx)
	case f.Config.Splunk != nil:
		return f.Config.Splunk.init(ctx)
	case f.Config.Syslog != nil:
		return f.Config.Syslog.init(ctx)
	case f.Config.Datadog != nil:
		return f.Config.Datadog.init(ctx)
	case f.Config.Amqp != nil:
		return f.Config.Amqp.init(ctx)
	case f.Config.Otlp != nil:
		return f.Config.Otlp.init(ctx)
	case f.Config.Elastic != nil:
		return f.Config.Elastic.init(ctx)
	case f.Config.Clickhouse != nil:
		return f.Config.Clickhouse.init(ctx)
	case f.Config.CloudWatch != nil:
		return f.Config.CloudWatch.init(ctx)
	case f.Config.GoogleCloudLogging != nil:
		return f.Config.GoogleCloudLogging.init(ctx)
	default:
		return fmt.Errorf("unsupported forwarder type")
	}
}

// Close releases whatever the active backend acquired in Init.
func (f *ForwarderV2) Close(ctx context.Context) error {
	switch {
	case f.Config.Http != nil:
		return f.Config.Http.close(ctx)
	case f.Config.Splunk != nil:
		return f.Config.Splunk.close(ctx)
	case f.Config.Syslog != nil:
		return f.Config.Syslog.close(ctx)
	case f.Config.Datadog != nil:
		return f.Config.Datadog.close(ctx)
	case f.Config.Amqp != nil:
		return f.Config.Amqp.close(ctx)
	case f.Config.Otlp != nil:
		return f.Config.Otlp.close(ctx)
	case f.Config.Elastic != nil:
		return f.Config.Elastic.close(ctx)
	case f.Config.Clickhouse != nil:
		return f.Config.Clickhouse.close(ctx)
	case f.Config.CloudWatch != nil:
		return f.Config.CloudWatch.close(ctx)
	case f.Config.GoogleCloudLogging != nil:
		return f.Config.GoogleCloudLogging.close(ctx)
	default:
		return fmt.Errorf("unsupported forwarder type")
	}
}

// Call forwards a single record to the active backend.
func (f *ForwarderV2) Call(ctx context.Context, record *LogRecord) error {
	switch {
	case f.Config.Http != nil:
		return f.Config.Http.call(ctx, record)
	case f.Config.Splunk != nil:
		return f.Config.Splunk.call(ctx, record)
	case f.Config.Syslog != nil:
		return f.Config.Syslog.call(ctx, record)
	case f.Config.Datadog != nil:
		return f.Config.Datadog.call(ctx, record)
	case f.Config.Amqp != nil:
		return f.Config.Amqp.call(ctx, record)
	case f.Config.Otlp != nil:
		return f.Config.Otlp.call(ctx, record)
	case f.Config.Elastic != nil:
		return f.Config.Elastic.call(ctx, record)
	case f.Config.Clickhouse != nil:
		return f.Config.Clickhouse.call(ctx, record)
	case f.Config.CloudWatch != nil:
		return f.Config.CloudWatch.call(ctx, record)
	case f.Config.GoogleCloudLogging != nil:
		return f.Config.GoogleCloudLogging.call(ctx, record)
	default:
		return fmt.Errorf("unsupported forwarder type")
	}
}

// MarshalJSON serialises the union as the JSON of whichever backend is set.
func (cfg *ForwarderConfigV2) MarshalJSON() ([]byte, error) {
	switch {
	case cfg.Http != nil:
		return json.Marshal(&cfg.Http)

	case cfg.Syslog != nil:
		return json.Marshal(&cfg.Syslog)

	case cfg.Datadog != nil:
		return json.Marshal(&cfg.Datadog)

	case cfg.Amqp != nil:
		return json.Marshal(&cfg.Amqp)

	case cfg.Splunk != nil:
		return json.Marshal(&cfg.Splunk)

	case cfg.Otlp != nil:
		return json.Marshal(&cfg.Otlp)

	case cfg.Elastic != nil:
		return json.Marshal(&cfg.Elastic)

	case cfg.Clickhouse != nil:
		return json.Marshal(&cfg.Clickhouse)

	case cfg.CloudWatch != nil:
		return json.Marshal(&cfg.CloudWatch)

	case cfg.GoogleCloudLogging != nil:
		return json.Marshal(&cfg.GoogleCloudLogging)

	default:
		return nil, fmt.Errorf("unsupported forwarder type")
	}
}

// UnmarshalJSON resets the union and decodes into the backend selected by the
// payload's "type" discriminator.
func (cfg *ForwarderConfigV2) UnmarshalJSON(data []byte) error {
	cfg.Http = nil
	cfg.Syslog = nil
	cfg.Datadog = nil
	cfg.Amqp = nil
	cfg.Splunk = nil
	cfg.Otlp = nil
	cfg.Elastic = nil
	cfg.Clickhouse = nil
	cfg.CloudWatch = nil
	cfg.GoogleCloudLogging = nil

	var typeInfo struct {
		Type string `json:"type" required:"true"`
	}

	if err := json.Unmarshal(data, &typeInfo); err != nil {
		return fmt.Errorf("failed to unmarshal forwarder type: %w", err)
	}

	switch typeInfo.Type {
	case "http":
		return json.Unmarshal(data, &cfg.Http)

	case "syslog":
		return json.Unmarshal(data, &cfg.Syslog)

	case "datadog":
		return json.Unmarshal(data, &cfg.Datadog)

	case "amqp":
		return json.Unmarshal(data, &cfg.Amqp)

	case "splunk":
		return json.Unmarshal(data, &cfg.Splunk)

	case "otlp":
		return json.Unmarshal(data, &cfg.Otlp)

	case "elastic":
		return json.Unmarshal(data, &cfg.Elastic)

	case "clickhouse":
		return json.Unmarshal(data, &cfg.Clickhouse)

	case "cloudwatch":
		return json.Unmarshal(data, &cfg.CloudWatch)

	case "googlecloudlogging":
		return json.Unmarshal(data, &cfg.GoogleCloudLogging)

	default:
		return fmt.Errorf("unsupported forwarder type: %s", typeInfo.Type)
	}
}
