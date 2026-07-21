package models

import (
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
	AwsCloudWatch      *ForwarderAwsCloudWatchV2      `json:"-"`
	GoogleCloudLogging *ForwarderGoogleCloudLoggingV2 `json:"-"`
	AzureMonitor       *ForwarderAzureMonitorV2       `json:"-"`
}

// JSONSchemaOneOf advertises every backend variant so the generated OpenAPI
// schema models Config as a "oneOf".
func (*ForwarderConfigV2) JSONSchemaOneOf() []any {
	return []any{
		ForwarderHttpV2{},
		ForwarderSyslogV2{},
		ForwarderDatadogV2{},
		ForwarderAmqpV2{},
		ForwarderSplunkV2{},
		ForwarderOtlpV2{},
		ForwarderElasticV2{},
		ForwarderClickhouseV2{},
		ForwarderAwsCloudWatchV2{},
		ForwarderGoogleCloudLoggingV2{},
		ForwarderAzureMonitorV2{},
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

	case cfg.AwsCloudWatch != nil:
		return json.Marshal(&cfg.AwsCloudWatch)

	case cfg.GoogleCloudLogging != nil:
		return json.Marshal(&cfg.GoogleCloudLogging)

	case cfg.AzureMonitor != nil:
		return json.Marshal(&cfg.AzureMonitor)

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
	cfg.AwsCloudWatch = nil
	cfg.GoogleCloudLogging = nil
	cfg.AzureMonitor = nil

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

	case "awscloudwatch":
		return json.Unmarshal(data, &cfg.AwsCloudWatch)

	case "googlecloudlogging":
		return json.Unmarshal(data, &cfg.GoogleCloudLogging)

	case "azuremonitor":
		return json.Unmarshal(data, &cfg.AzureMonitor)

	default:
		return fmt.Errorf("unsupported forwarder type: %s", typeInfo.Type)
	}
}
