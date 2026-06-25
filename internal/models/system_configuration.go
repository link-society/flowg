package models

// SystemConfiguration holds the global, server-wide settings. SyslogAllowedOrigins
// restricts which source IPs or CIDR ranges may push logs to the syslog endpoint.
type SystemConfiguration struct {
	SyslogAllowedOrigins []string `json:"syslog_allowed_origins"`
}
