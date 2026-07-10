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
// exists through their "index:user:<username>" marker, then stores the token
// hash under a new "pat:<username>:<uuid>" key and returns the plaintext token
// (shown to the caller this one time only) along with its UUID.
func CreateToken(txn kv.MutationTx, username string) (string, string, error) {
	token, err := secret.NewSecret("pat", 32)
	if err != nil {
		return "", "", err
	}

	tokenHash, err := hash.HashPassword(token)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash token for user %q: %w", username, err)
	}

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
	err = txn.Set(tokenKey, []byte(tokenHash))
	if err != nil {
		return "", "", fmt.Errorf("failed to write token for user %q: %w", username, err)
	}

	return token, tokenUuid, nil
}

// VerifyToken resolves a plaintext token back to its owning user. Tokens are
// not addressable by their value, so it scans every "pat:" key and bcrypt-
// compares the candidate against each stored hash; on a match it returns the
// owning user (the username is the key segment between the "pat:" prefix and the
// trailing UUID).
func VerifyToken(txn kv.QueryTx, token string) (*models.User, error) {
	var username string
	found := false

	for kv := range txn.IterPairs(kv.Key{"pat"}, kv.KeyRange{}) {
		key := kv.Key()
		val := kv.Value()

		associatedUser := key[len(key)-2]

		tokenHash := string(val)
		isValid, err := hash.VerifyPassword(token, tokenHash)
		if err != nil {
			return nil, fmt.Errorf("failed to verify token: %w", err)
		}

		if isValid {
			username = associatedUser
			found = true
			break
		}
	}

	if !found {
		return nil, nil
	}

	return FetchUser(txn, username)
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

// DeleteToken removes a single PAT, identified by its UUID, from a user.
func DeleteToken(txn kv.MutationTx, username string, tokenUUID string) error {
	key := kv.Key{"pat", username, tokenUUID}
	err := txn.Clear(key)
	if err != nil {
		return fmt.Errorf("failed to clear token %q for user %q: %w", tokenUUID, username, err)
	}

	return nil
}
