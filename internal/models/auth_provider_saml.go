package models

// AuthProviderSaml contains the configuration for a SAML authentication
// provider.
type AuthProviderSaml struct {
	Type           string `json:"type" enum:"saml" required:"true"`
	IdpMetadataURL string `json:"idp_metadata_url" required:"true" format:"uri"`
}