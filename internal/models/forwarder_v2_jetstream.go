package models

// ForwarderJetstreamV2 forwards records to an nats.io Jetstream,
type ForwarderJetstreamV2 struct {
	Type           string                           `json:"type" enum:"jetstream" required:"true"`
	Servers        []string                         `json:"servers"`
	Subject        ForwarderJetstreamSubjectField   `json:"subject" required:"true"`
	Body           ForwarderJetstreamBodyField      `json:"body" required:"true"`
	MessageID      ForwarderJetstreamMessageIDField `json:"message_id"`
	ExpectedStream string                           `json:"expected_stream"`
	PublishTimeout int                              `json:"publish_timeout"`
	Headers        map[string][]string              `json:"headers"`

	AuthMode    string `json:"auth_mode"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	Credentials string `json:"credentials"`

	TlsCA                 string `json:"tls_ca"`
	TlsCertificate        string `json:"tls_certificate"`
	TlsPrivateKey         string `json:"tls_private_key"`
	TlsServerName         string `json:"tls_server_name"`
	TlsInsecureSkipVerify bool   `json:"tls_insecure_skip_verify"`
}

type ForwarderJetstreamBodyField string
type ForwarderJetstreamSubjectField string
type ForwarderJetstreamMessageIDField string

func (ForwarderJetstreamBodyField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderJetstreamSubjectField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderJetstreamMessageIDField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}
