package models

// ForwarderDatadogV2 forwards records to the Datadog logs intake. The ddsource,
// ddtags, hostname, message and service attributes are dynamic fields evaluated
// per record.
type ForwarderDatadogV2 struct {
	Type     string                          `json:"type" enum:"datadog" required:"true"`
	Url      string                          `json:"url" required:"true" format:"uri"`
	ApiKey   string                          `json:"apiKey" required:"true" minLength:"1"`
	DDsource ForwarderDatadogV2DDsourceField `json:"ddsource" required:"true"`
	DDtags   ForwarderDatadogV2DDtagsField   `json:"ddtags" required:"true"`
	Hostname ForwarderDatadogV2HostnameField `json:"hostname" required:"true"`
	Message  ForwarderDatadogV2MessageField  `json:"message" required:"true"`
	Service  ForwarderDatadogV2ServiceField  `json:"service" required:"true"`
}

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
