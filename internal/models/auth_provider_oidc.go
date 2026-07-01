package models

// AuthProviderOidc contains the configuration for an OpenID Connect
// authentication provider.
type AuthProviderOidc struct {
	Type         string `json:"type" enum:"oidc" required:"true"`
	Issuer       string `json:"issuer" required:"true" format:"uri"`
	ClientID     string `json:"client_id" required:"true"`
	ClientSecret string `json:"client_secret" required:"true"`
}