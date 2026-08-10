package transactions

import (
	"encoding/json"
	"fmt"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"
)

const PROVIDER = "provider"

// ListAuthProviders returns deserialized array of models.AuthProvider
func ListAuthProviders(txn kv.QueryTx) ([]models.AuthProvider, error) {
	var providers []models.AuthProvider

	for key := range txn.IterKeys(kv.Key{PROVIDER}, kv.KeyRange{}) {
		var provider models.AuthProvider
		err := json.Unmarshal([]byte(key[len(key)-1]), &provider)

		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall auth provider: %w", err)
		}

		providers = append(providers, provider)
	}

	return providers, nil
}

// ReadAuthProvider returns deserialized instance of models.AuthProvider with provided type and name
func ReadAuthProvider(txn kv.QueryTx, providerType string, name string) (models.AuthProvider, error) {
	key := kv.Key{PROVIDER, providerType, name}

	marshalled, err := txn.Get(key)
	if err != nil {
		return models.AuthProvider{}, fmt.Errorf("failed to find auth provider %q with type %q: %w", name, providerType, err)
	}

	var provider models.AuthProvider
	err = json.Unmarshal(marshalled, &provider)
	if err != nil {
		return models.AuthProvider{}, fmt.Errorf("failed to unmarshall auth provider %q with type %q: %w", name, providerType, err)
	}

	return provider, nil
}

// SaveAuthProvider serializes models.AuthProvider and saves it to the database
func SaveAuthProvider(txn kv.MutationTx, provider models.AuthProvider) error {
	var providerType string

	if provider.Config.Oidc != nil {
		providerType = provider.Config.Oidc.Type
	}

	if provider.Config.Saml != nil {
		providerType = provider.Config.Saml.Type
	}

	key := kv.Key{PROVIDER, providerType, provider.Name}

	marshalled, err := json.Marshal(provider)
	if err != nil {
		return fmt.Errorf("failed to marshall auth provider %q with type %q: %w", provider.Name, providerType, err)
	}

	return txn.Set(key, marshalled)
}

// DeleteAuthProvider deletes provider from the database
func DeleteAuthProvider(txn kv.MutationTx, providerType string, name string) error {
	providerKey := kv.Key{PROVIDER, providerType, name}

	if err := txn.Clear(providerKey); err != nil {
		return fmt.Errorf("failed to clear auth provider %q: %w", name, err)
	}

	return nil
}
