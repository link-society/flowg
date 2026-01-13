package models

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
