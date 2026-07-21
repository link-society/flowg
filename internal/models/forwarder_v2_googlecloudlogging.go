package models

// ForwarderGoogleCloudLoggingV2 forwards records to Google Cloud Watch Logs stream,
// authenticating with JSON credentials if available.
type ForwarderGoogleCloudLoggingV2 struct {
	Type        string `json:"type" enum:"googlecloudlogging" required:"true"`
	Endpoint    string `json:"endpoint" required:"true"`
	ProjectID   string `json:"project_id" required:"true"`
	LogID       string `json:"log_id" required:"true"`
	DisableTLS  bool   `json:"disable_tls"`
	DisableAuth bool   `json:"disable_auth"`
	AuthJSON    string `json:"auth_json"`
}
