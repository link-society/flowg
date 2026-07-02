package models_test

import (
	"reflect"
	"testing"

	"encoding/json"

	"link-society.com/flowg/internal/models"
)

func TestAuthProvider_RoundTrip_Oidc(t *testing.T) {
	original := models.AuthProvider{
		Name:        "my-oidc",
		DisplayName: "My OIDC Provider",
		Config: models.AuthProviderConfig{
			Oidc: &models.AuthProviderOidc{
				Type:         "oidc",
				Issuer:       "https://accounts.example.com",
				ClientID:     "client-123",
				ClientSecret: "secret-456",
			},
		},
	}

	data, err := json.Marshal(&original)
	if err != nil {
		t.Fatalf("failed to marshal AuthProvider: %v", err)
	}

	var decoded models.AuthProvider
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal AuthProvider: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Fatalf("round-trip mismatch:\n  original = %+v\n  decoded  = %+v", original, decoded)
	}
}

func TestAuthProvider_RoundTrip_Saml(t *testing.T) {
	original := models.AuthProvider{
		Name:        "my-saml",
		DisplayName: "My SAML Provider",
		Config: models.AuthProviderConfig{
			Saml: &models.AuthProviderSaml{
				Type:           "saml",
				IdpMetadataURL: "https://idp.example.com/metadata",
			},
		},
	}

	data, err := json.Marshal(&original)
	if err != nil {
		t.Fatalf("failed to marshal AuthProvider: %v", err)
	}

	var decoded models.AuthProvider
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal AuthProvider: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Fatalf("round-trip mismatch:\n  original = %+v\n  decoded  = %+v", original, decoded)
	}
}

func TestAuthProviderConfig_Marshal_NoVariant(t *testing.T) {
	provider := models.AuthProvider{Name: "broken", DisplayName: "No variant"}
	if _, err := json.Marshal(&provider); err == nil {
		t.Fatal("expected an error marshalling an AuthProvider with no variant, got nil")
	}
}
