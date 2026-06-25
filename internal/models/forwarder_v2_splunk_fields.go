package models

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
