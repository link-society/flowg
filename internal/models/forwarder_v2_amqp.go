package models

// ForwarderAmqpV2 publishes records to an AMQP (e.g. RabbitMQ) exchange. The
// exchange, routing key and body are dynamic fields evaluated per record.
type ForwarderAmqpV2 struct {
	Type       string                         `json:"type" enum:"amqp" required:"true"`
	Url        string                         `json:"url" required:"true" format:"uri"`
	Exchange   ForwarderAmqpV2ExchangeField   `json:"exchange" required:"true" minLength:"1"`
	RoutingKey ForwarderAmqpV2RoutingKeyField `json:"routing_key" default:""`
	Body       ForwarderAmqpV2BodyField       `json:"body,omitempty"`
}

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
