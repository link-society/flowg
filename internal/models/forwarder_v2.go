package models

import (
	"context"
	"fmt"

	"encoding/json"
)

type ForwarderV2 struct {
	Version int                `json:"version" default:"2"`
	Config  *ForwarderConfigV2 `json:"config"`
}

type ForwarderConfigV2 struct {
	Http   *ForwarderHttpV2   `json:"-"`
	Syslog *ForwarderSyslogV2 `json:"-"`
}

func (ForwarderConfigV2) JSONSchemaOneOf() []any {
	return []any{
		ForwarderHttpV2{},
		ForwarderSyslogV2{},
	}
}

func (f *ForwarderV2) Call(ctx context.Context, record *LogRecord) error {
	switch {
	case f.Config.Http != nil:
		return f.Config.Http.call(ctx, record)

	case f.Config.Syslog != nil:
		return f.Config.Syslog.call(ctx, record)

	default:
		return fmt.Errorf("unsupported forwarder type")
	}
}

func (cfg *ForwarderConfigV2) MarshalJSON() ([]byte, error) {
	switch {
	case cfg.Http != nil:
		return json.Marshal(&cfg.Http)

	case cfg.Syslog != nil:
		return json.Marshal(&cfg.Syslog)

	default:
		return nil, fmt.Errorf("unsupported forwarder type")
	}
}

func (cfg *ForwarderConfigV2) UnmarshalJSON(data []byte) error {
	var typeInfo struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &typeInfo); err != nil {
		return fmt.Errorf("failed to unmarshal forwarder type: %w", err)
	}

	switch typeInfo.Type {
	case "http":
		return json.Unmarshal(data, &cfg.Http)

	case "syslog":
		return json.Unmarshal(data, &cfg.Syslog)

	default:
		return fmt.Errorf("unsupported forwarder type: %s", typeInfo.Type)
	}
}
