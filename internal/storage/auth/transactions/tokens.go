package transactions

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"

	"link-society.com/flowg/internal/models"

	authUtils "link-society.com/flowg/internal/utils/auth"
	"link-society.com/flowg/internal/utils/auth/hash"
)

func CreateToken(txn *badger.Txn, username string) (string, string, error) {
	token, err := authUtils.NewSecret("pat", 32)
	if err != nil {
		return "", "", err
	}

	tokenHash, err := hash.HashPassword(token)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash token: %w", err)
	}

	tokenUuid := uuid.New().String()

	userKey := []byte(fmt.Sprintf("index:user:%s", username))
	_, err = txn.Get(userKey)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return "", "", fmt.Errorf("user '%s' does not exist", username)
		}

		return "", "", fmt.Errorf("failed to check if user '%s' exists: %w", username, err)
	}

	tokenKey := []byte(fmt.Sprintf("pat:%s:%s", username, tokenUuid))
	err = txn.Set(tokenKey, []byte(tokenHash))
	if err != nil {
		return "", "", fmt.Errorf("failed to add token to user '%s': %w", username, err)
	}

	return token, tokenUuid, nil
}

func VerifyToken(txn *badger.Txn, token string) (*models.User, error) {
	var username string
	found := false

	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte("pat:")
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().Key()
		keySuffix := string(key[4:])
		associatedUser := keySuffix[:len(keySuffix)-37] // remove UUID

		err := it.Item().Value(func(val []byte) error {
			tokenHash := string(val)
			isValid, err := hash.VerifyPassword(token, tokenHash)
			if err != nil {
				return fmt.Errorf("failed to verify token: %w", err)
			}

			if isValid {
				username = associatedUser
				found = true
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		if found {
			break
		}
	}

	if !found {
		return nil, nil
	}

	return FetchUser(txn, username)
}

func ListTokens(txn *badger.Txn, username string) []string {
	tokenUUIDs := []string{}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte(fmt.Sprintf("pat:%s:", username))
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().Key()
		tokenUUIDs = append(tokenUUIDs, string(key[len(username)+5:]))
	}

	return tokenUUIDs
}

func DeleteToken(txn *badger.Txn, username string, tokenUUID string) error {
	key := []byte(fmt.Sprintf("pat:%s:%s", username, tokenUUID))
	err := txn.Delete(key)
	if err != nil && err != badger.ErrKeyNotFound {
		return fmt.Errorf(
			"failed to delete token '%s' for user '%s': %w",
			tokenUUID, username, err,
		)
	}

	return nil
}
