package models

// The Forwarder*Field types below are the AMQP forwarder's per-record fields:
// each is either a literal string or a DynamicField ("@expr:" expression), as
// advertised by their JSONSchemaAnyOf methods.
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
