package models

import (
	"encoding/json"
	"fmt"
)

type AuthProvider struct {
	Name        string             `json:"name" required:"true"`
	DisplayName string             `json:"display_name" required:"true"`
	Config      AuthProviderConfig `json:"config" required:"true"`
}

type AuthProviderConfig struct {
	Oidc *OidcAuthProvider `json:"-"`
	Saml *SamlAuthProvider `json:"-"`
}

// JSONSchemaOneOf advertises every provider variant so the generated OpenAPI
// schema models Config as a "oneOf".
func (AuthProviderConfig) JSONSchemaOneOf() []any {
	return []any{
		OidcAuthProvider{},
		SamlAuthProvider{},
	}
}

func (cfg *AuthProviderConfig) MarshalJSON() ([]byte, error) {
	switch {
	case cfg.Oidc != nil:
		return json.Marshal(&cfg.Oidc)

	case cfg.Saml != nil:
		return json.Marshal(&cfg.Saml)

	default:
		return nil, fmt.Errorf("unsupported auth provider type")
	}
}

func (cfg *AuthProviderConfig) UnmarshalJSON(data []byte) error {
	cfg.Oidc = nil
	cfg.Saml = nil

	var typeInfo struct {
		Type string `json:"type" required:"true"`
	}

	if err := json.Unmarshal(data, &typeInfo); err != nil {
		return fmt.Errorf("failed to unmarshal auth provider type: %w", err)
	}

	switch typeInfo.Type {
	case "oidc":
		return json.Unmarshal(data, &cfg.Oidc)

	case "saml":
		return json.Unmarshal(data, &cfg.Saml)

	default:
		return fmt.Errorf("unsupported auth provider type: %s", typeInfo.Type)
	}
}