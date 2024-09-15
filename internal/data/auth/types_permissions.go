package auth

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

type Permissions struct {
	CanViewPipelines bool `json:"can_view_pipelines"`
	CanEditPipelines bool `json:"can_edit_pipelines"`

	CanViewTransformers bool `json:"can_view_transformers"`
	CanEditTransformers bool `json:"can_edit_transformers"`

	CanViewStreams bool `json:"can_view_streams"`
	CanEditStreams bool `json:"can_edit_streams"`

	CanViewAlerts bool `json:"can_view_alerts"`
	CanEditAlerts bool `json:"can_edit_alerts"`

	CanViewACLs bool `json:"can_view_acls"`
	CanEditACLs bool `json:"can_edit_acls"`

	CanSendLogs bool `json:"can_send_logs"`
}

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

func PermissionsFromScopes(scopes []Scope) Permissions {
	permissions := Permissions{}
	for _, scope := range scopes {
		switch scope {
		case SCOPE_READ_PIPELINES:
			permissions.CanViewPipelines = true
		case SCOPE_WRITE_PIPELINES:
			permissions.CanEditPipelines = true
		case SCOPE_READ_TRANSFORMERS:
			permissions.CanViewTransformers = true
		case SCOPE_WRITE_TRANSFORMERS:
			permissions.CanEditTransformers = true
		case SCOPE_READ_STREAMS:
			permissions.CanViewStreams = true
		case SCOPE_WRITE_STREAMS:
			permissions.CanEditStreams = true
		case SCOPE_READ_ALERTS:
			permissions.CanViewAlerts = true
		case SCOPE_WRITE_ALERTS:
			permissions.CanEditAlerts = true
		case SCOPE_READ_ACLS:
			permissions.CanViewACLs = true
		case SCOPE_WRITE_ACLS:
			permissions.CanEditACLs = true
		case SCOPE_SEND_LOGS:
			permissions.CanSendLogs = true
		}
	}
	return permissions
}
