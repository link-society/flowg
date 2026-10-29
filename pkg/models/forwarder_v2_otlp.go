package models

// ForwarderOtlpV2 forwards records to an OpenTelemetry (OTLP/HTTP) logs
// endpoint, encoding them as protobuf.
type ForwarderOtlpV2 struct {
	Type     string            `json:"type" enum:"otlp" required:"true"`
	Endpoint string            `json:"endpoint" required:"true" format:"uri"`
	Headers  map[string]string `json:"headers,omitempty"`
}
