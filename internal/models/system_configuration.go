package models

type SystemConfiguration struct {
	SyslogAllowedOrigins []string `json:"syslog_allowed_origins"`
}
