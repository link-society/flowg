package models

type Permissions struct {
	CanViewPipelines bool `json:"can_view_pipelines" required:"true"`
	CanEditPipelines bool `json:"can_edit_pipelines" required:"true"`

	CanViewTransformers bool `json:"can_view_transformers" required:"true"`
	CanEditTransformers bool `json:"can_edit_transformers" required:"true"`

	CanViewStreams bool `json:"can_view_streams" required:"true"`
	CanEditStreams bool `json:"can_edit_streams" required:"true"`

	CanViewForwarders bool `json:"can_view_forwarders" required:"true"`
	CanEditForwarders bool `json:"can_edit_forwarders" required:"true"`

	CanViewACLs bool `json:"can_view_acls" required:"true"`
	CanEditACLs bool `json:"can_edit_acls" required:"true"`

	CanSendLogs bool `json:"can_send_logs" required:"true"`
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
		case SCOPE_READ_FORWARDERS:
			permissions.CanViewForwarders = true
		case SCOPE_WRITE_FORWARDERS:
			permissions.CanEditForwarders = true
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
