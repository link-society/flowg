package models

type OidcAuthProvider struct {
	Type         string `json:"type" enum:"oidc" required:"true"`
	Issuer       string `json:"issuer" required:"true" format:"uri"`
	ClientID     string `json:"client_id" required:"true"`
	ClientSecret string `json:"client_secret" required:"true"`
}