package models

// ForwarderElasticV2 indexes records into an Elasticsearch index, optionally
// trusting a custom CA certificate.
type ForwarderElasticV2 struct {
	Type      string   `json:"type" enum:"elastic" required:"true"`
	Index     string   `json:"index" required:"true" minLength:"1"`
	Username  string   `json:"username" required:"true" minLength:"1"`
	Password  string   `json:"password" required:"true" minLength:"1"`
	Addresses []string `json:"addresses" required:"true" minItems:"1" items.format:"uri"`
	CACert    string   `json:"ca,omitempty" format:"pem"`
}
