package transactions

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"

	"link-society.com/flowg/internal/models"

	authUtils "link-society.com/flowg/internal/utils/auth"
	"link-society.com/flowg/internal/utils/auth/hash"
	"link-society.com/flowg/internal/utils/hlc"
)

func CreateToken(txn *badger.Txn, username string, ts hlc.Timestamp) (string, string, error) {
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
	_, found, err := getItem(txn, userKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to check if user '%s' exists: %w", username, err)
	}
	if !found {
		return "", "", fmt.Errorf("user '%s' does not exist", username)
	}

	tokenKey := []byte(fmt.Sprintf("pat:%s:%s", username, tokenUuid))
	if err := setItem(txn, tokenKey, []byte(tokenHash), ts); err != nil {
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
		item := it.Item()
		payload, live, err := liveValue(item)
		if err != nil {
			return nil, err
		}
		if !live {
			continue
		}

		keySuffix := string(item.Key()[4:])
		associatedUser := keySuffix[:len(keySuffix)-37] // remove UUID

		isValid, err := hash.VerifyPassword(token, string(payload))
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

func ListTokens(txn *badger.Txn, username string) ([]string, error) {
	tokenUUIDs := []string{}

	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte(fmt.Sprintf("pat:%s:", username))
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		_, live, err := liveValue(item)
		if err != nil {
			return nil, err
		}
		if !live {
			continue
		}
		tokenUUIDs = append(tokenUUIDs, string(item.Key()[len(username)+5:]))
	}

	return tokenUUIDs, nil
}

func DeleteToken(txn *badger.Txn, username string, tokenUUID string, ts hlc.Timestamp) error {
	key := []byte(fmt.Sprintf("pat:%s:%s", username, tokenUUID))
	if err := deleteItem(txn, key, ts); err != nil {
		return fmt.Errorf(
			"failed to delete token '%s' for user '%s': %w",
			tokenUUID, username, err,
		)
	}

	return nil
}
