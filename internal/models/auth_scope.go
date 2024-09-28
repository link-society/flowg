package models

import "fmt"

type Scope string

const (
	SCOPE_READ_PIPELINES     Scope = "read_pipelines"
	SCOPE_WRITE_PIPELINES    Scope = "write_pipelines"
	SCOPE_READ_TRANSFORMERS  Scope = "read_transformers"
	SCOPE_WRITE_TRANSFORMERS Scope = "write_transformers"
	SCOPE_READ_STREAMS       Scope = "read_streams"
	SCOPE_WRITE_STREAMS      Scope = "write_streams"
	SCOPE_READ_ALERTS        Scope = "read_alerts"
	SCOPE_WRITE_ALERTS       Scope = "write_alerts"
	SCOPE_READ_ACLS          Scope = "read_acls"
	SCOPE_WRITE_ACLS         Scope = "write_acls"
	SCOPE_SEND_LOGS          Scope = "send_logs"
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
	case "read_alerts":
		return SCOPE_READ_ALERTS, nil
	case "write_alerts":
		return SCOPE_WRITE_ALERTS, nil
	case "read_acls":
		return SCOPE_READ_ACLS, nil
	case "write_acls":
		return SCOPE_WRITE_ACLS, nil
	case "send_logs":
		return SCOPE_SEND_LOGS, nil
	default:
		return "", fmt.Errorf("invalid scope: %s", s)
	}
}
