package models

import (
	"encoding/json"
	"fmt"
)

// AuthProvider is the current authentication provider model for delegated
// authentication. The concrete provider lives in Config, which is a tagged
// union of one backend type.
type AuthProvider struct {
	Name        string             `json:"name" required:"true"`
	DisplayName string             `json:"display_name" required:"true"`
	Config      AuthProviderConfig `json:"config" required:"true"`
}

// AuthProviderConfig is a tagged union: exactly one field is non-nil, selecting
// the authentication provider backend.
type AuthProviderConfig struct {
	Oidc *AuthProviderOidc `json:"-"`
	Saml *AuthProviderSaml `json:"-"`
}

// JSONSchemaOneOf advertises every provider variant so the generated OpenAPI
// schema models Config as a "oneOf".
func (AuthProviderConfig) JSONSchemaOneOf() []any {
	return []any{
		AuthProviderOidc{},
		AuthProviderSaml{},
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