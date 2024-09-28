package models

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
