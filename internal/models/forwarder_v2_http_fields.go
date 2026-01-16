package models

type ForwarderHttpV2BodyField string

func (ForwarderHttpV2BodyField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}
