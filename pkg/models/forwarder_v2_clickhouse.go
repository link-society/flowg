package models

// ForwarderClickhouseV2 inserts records into a ClickHouse table, creating it on
// first use. Each record becomes a row of (id, timestamp, fields map).
type ForwarderClickhouseV2 struct {
	Type     string `json:"type" enum:"clickhouse" required:"true"`
	Address  string `json:"address" required:"true" pattern:"^(([a-zA-Z0-9.-]+)|(\\[[0-9A-Fa-f:]+\\])):[0-9]{1,5}$"`
	Database string `json:"db" required:"true" minLength:"1"`
	Table    string `json:"table" required:"true" pattern:"^[a-zA-Z_][a-zA-Z0-9_]*$" minLength:"1" maxLength:"64"`
	Username string `json:"user" required:"true" minLength:"1"`
	Password string `json:"pass" required:"true" minLength:"1"`
	UseTls   bool   `json:"tls" required:"true"`
}
