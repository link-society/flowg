package models

type ForwarderAmqpV2ExchangeField string
type ForwarderAmqpV2RoutingKeyField string
type ForwarderAmqpV2BodyField string

func (ForwarderAmqpV2ExchangeField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderAmqpV2RoutingKeyField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderAmqpV2BodyField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}
