package models

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
