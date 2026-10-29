package models

// ForwarderSplunkV2 forwards records to a Splunk HTTP Event Collector. The
// source and host fields are dynamic fields evaluated per record.
type ForwarderSplunkV2 struct {
	Type     string                       `json:"type" enum:"splunk" required:"true"`
	Endpoint string                       `json:"endpoint" required:"true" format:"uri"`
	Token    string                       `json:"token" required:"true" minLength:"1"`
	Source   ForwarderSplunkV2SourceField `json:"source"`
	Host     ForwarderSplunkV2HostField   `json:"host"`
}

// The Forwarder*Field types below are the Splunk forwarder's per-record fields:
// each is either a literal string or a DynamicField ("@expr:" expression), as
// advertised by their JSONSchemaAnyOf methods.
type ForwarderSplunkV2SourceField string
type ForwarderSplunkV2HostField string

func (ForwarderSplunkV2SourceField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderSplunkV2HostField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}
