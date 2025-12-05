package models

import "fmt"

type Scope string

const (
	SCOPE_READ_PIPELINES             Scope = "read_pipelines"
	SCOPE_WRITE_PIPELINES            Scope = "write_pipelines"
	SCOPE_READ_TRANSFORMERS          Scope = "read_transformers"
	SCOPE_WRITE_TRANSFORMERS         Scope = "write_transformers"
	SCOPE_READ_STREAMS               Scope = "read_streams"
	SCOPE_WRITE_STREAMS              Scope = "write_streams"
	SCOPE_READ_FORWARDERS            Scope = "read_forwarders"
	SCOPE_WRITE_FORWARDERS           Scope = "write_forwarders"
	SCOPE_READ_ACLS                  Scope = "read_acls"
	SCOPE_WRITE_ACLS                 Scope = "write_acls"
	SCOPE_SEND_LOGS                  Scope = "send_logs"
	SCOPE_READ_SYSTEM_CONFIGURATION  Scope = "read_system_configuration"
	SCOPE_WRITE_SYSTEM_CONFIGURATION Scope = "write_system_configuration"
)

func ParseScope(s string) (Scope, error) {
	switch s {
	case "read_pipelines":
		return SCOPE_READ_PIPELINES, nil
	case "write_pipelines":
		return SCOPE_WRITE_PIPELINES, nil
	case "read_transformers":
		return SCOPE_READ_TRANSFORMERS, nil
	case "write_transformers":
		return SCOPE_WRITE_TRANSFORMERS, nil
	case "read_streams":
		return SCOPE_READ_STREAMS, nil
	case "write_streams":
		return SCOPE_WRITE_STREAMS, nil
	case "read_forwarders":
		return SCOPE_READ_FORWARDERS, nil
	case "write_forwarders":
		return SCOPE_WRITE_FORWARDERS, nil
	case "read_acls":
		return SCOPE_READ_ACLS, nil
	case "write_acls":
		return SCOPE_WRITE_ACLS, nil
	case "send_logs":
		return SCOPE_SEND_LOGS, nil
	case "read_system_configuration":
		return SCOPE_READ_SYSTEM_CONFIGURATION, nil
	case "write_system_configuration":
		return SCOPE_WRITE_SYSTEM_CONFIGURATION, nil
	default:
		return "", fmt.Errorf("invalid scope: %s", s)
	}
}

func (s Scope) Enum() []any {
	return []any{
		SCOPE_READ_PIPELINES,
		SCOPE_WRITE_PIPELINES,
		SCOPE_READ_TRANSFORMERS,
		SCOPE_WRITE_TRANSFORMERS,
		SCOPE_READ_STREAMS,
		SCOPE_WRITE_STREAMS,
		SCOPE_READ_FORWARDERS,
		SCOPE_WRITE_FORWARDERS,
		SCOPE_READ_ACLS,
		SCOPE_WRITE_ACLS,
		SCOPE_SEND_LOGS,
		SCOPE_READ_SYSTEM_CONFIGURATION,
		SCOPE_WRITE_SYSTEM_CONFIGURATION,
	}
}
