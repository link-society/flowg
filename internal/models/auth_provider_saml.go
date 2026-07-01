package models

type SamlAuthProvider struct {
	Type           string `json:"type" enum:"saml" required:"true"`
	IdpMetadataURL string `json:"idp_metadata_url" required:"true" format:"uri"`
}