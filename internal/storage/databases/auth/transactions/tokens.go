package transactions

import (
	"fmt"

	"github.com/google/uuid"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"

	"link-society.com/flowg/internal/utils/hash"
	"link-society.com/flowg/internal/utils/secret"
)

// CreateToken issues a fresh PAT for an existing user. It confirms the user
// exists through their "index:user:<username>" marker, stores the token's
// SHA-256 hash under "pat:<username>:<uuid>", and records a reverse index
// "index:pat:<hash>" -> username so the token can be resolved in O(1). It
// returns the plaintext token (shown to the caller this one time only) along
// with its UUID.
func CreateToken(txn kv.MutationTx, username string) (string, string, error) {
	token, err := secret.NewSecret("pat", 32)
	if err != nil {
		return "", "", err
	}

	tokenHash := hash.HashToken(token)
	tokenUuid := uuid.New().String()

	userKey := kv.Key{"index", "user", username}
	val, err := txn.Get(userKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to read user %q: %w", username, err)
	}
	if val == nil {
		return "", "", fmt.Errorf("user %q not found", username)
	}

	tokenKey := kv.Key{"pat", username, tokenUuid}
	if err := txn.Set(tokenKey, []byte(tokenHash)); err != nil {
		return "", "", fmt.Errorf("failed to write token for user %q: %w", username, err)
	}

	indexKey := kv.Key{"index", "pat", tokenHash}
	if err := txn.Set(indexKey, []byte(username)); err != nil {
		return "", "", fmt.Errorf("failed to index token for user %q: %w", username, err)
	}

	return token, tokenUuid, nil
}

// VerifyToken resolves a plaintext token back to its owning user in O(1): it
// hashes the token and looks up "index:pat:<hash>". A missing entry means the
// token is unknown (and no fallback to a legacy per-token scan is attempted).
func VerifyToken(txn kv.QueryTx, token string) (*models.User, error) {
	indexKey := kv.Key{"index", "pat", hash.HashToken(token)}

	val, err := txn.Get(indexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	if val == nil {
		return nil, nil
	}

	return FetchUser(txn, string(val))
}

// ListTokens returns the UUIDs of every PAT owned by a user by scanning the
// "pat:<username>:" prefix. Only the keys are needed, so values are not fetched.
func ListTokens(txn kv.QueryTx, username string) []string {
	tokenUUIDs := []string{}

	for key := range txn.IterKeys(kv.Key{"pat", username}, kv.KeyRange{}) {
		tokenUUIDs = append(tokenUUIDs, key[len(key)-1])
	}

	return tokenUUIDs
}

// DeleteToken removes a single PAT, identified by its UUID, from a user, along
// with its "index:pat:<hash>" reverse-index entry.
func DeleteToken(txn kv.MutationTx, username string, tokenUUID string) error {
	tokenKey := kv.Key{"pat", username, tokenUUID}

	tokenHash, err := txn.Get(tokenKey)
	if err != nil {
		return fmt.Errorf("failed to read token %q for user %q: %w", tokenUUID, username, err)
	}

	if tokenHash != nil {
		indexKey := kv.Key{"index", "pat", string(tokenHash)}
		if err := txn.Clear(indexKey); err != nil {
			return fmt.Errorf("failed to clear token index %q for user %q: %w", tokenUUID, username, err)
		}
	}

	if err := txn.Clear(tokenKey); err != nil {
		return fmt.Errorf("failed to clear token %q for user %q: %w", tokenUUID, username, err)
	}

	return nil
}
