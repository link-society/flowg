package models

// ForwarderHttpV2BodyField is the HTTP forwarder's body: either a literal string
// or a DynamicField (an "@expr:" expression). JSONSchemaAnyOf advertises both
// forms to the generated OpenAPI schema.
type ForwarderHttpV2BodyField string

func (ForwarderHttpV2BodyField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}
