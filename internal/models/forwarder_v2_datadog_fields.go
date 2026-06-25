package models

// The Forwarder*Field types below are the Datadog forwarder's per-record
// attributes: each is either a literal string or a DynamicField ("@expr:"
// expression), as advertised by their JSONSchemaAnyOf methods.
type ForwarderDatadogV2DDsourceField string
type ForwarderDatadogV2DDtagsField string
type ForwarderDatadogV2HostnameField string
type ForwarderDatadogV2MessageField string
type ForwarderDatadogV2ServiceField string

func (ForwarderDatadogV2DDsourceField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderDatadogV2DDtagsField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderDatadogV2HostnameField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderDatadogV2MessageField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderDatadogV2ServiceField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}
