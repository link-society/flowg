package models

// ForwarderHttpV2 forwards records to an arbitrary HTTP endpoint via POST. The
// request body is produced by evaluating the Body dynamic field (defaulting to
// the record serialised as JSON).
type ForwarderHttpV2 struct {
	Type    string                   `json:"type" enum:"http" required:"true"`
	Url     string                   `json:"url" required:"true" format:"uri"`
	Headers map[string]string        `json:"headers,omitempty"`
	Proxy   string                   `json:"proxy,omitempty"`
	Body    ForwarderHttpV2BodyField `json:"body,omitempty"`
}

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
