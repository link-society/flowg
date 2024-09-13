package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"

	"link-society.com/flowg/internal/data/auth/hash"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type TokenSystem struct {
	backend *Database
}

func NewTokenSystem(backend *Database) *TokenSystem {
	return &TokenSystem{backend: backend}
}

func (sys *TokenSystem) CreateToken(username string) (string, string, error) {
	token, err := newToken(32)
	if err != nil {
		return "", "", err
	}

	tokenHash, err := hash.HashPassword(token)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash token: %w", err)
	}

	tokenUuid := uuid.New().String()

	err = sys.backend.db.Update(func(txn *badger.Txn) error {
		userKey := []byte(fmt.Sprintf("index:user:%s", username))
		_, err := txn.Get(userKey)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return fmt.Errorf("user '%s' does not exist", username)
			}

			return fmt.Errorf("failed to check if user '%s' exists: %w", username, err)
		}

		tokenKey := []byte(fmt.Sprintf("pat:%s:%s", username, tokenUuid))
		err = txn.Set(tokenKey, []byte(tokenHash))
		if err != nil {
			return fmt.Errorf("failed to add token to user '%s': %w", username, err)
		}

		return nil
	})

	if err != nil {
		return "", "", err
	}

	return token, tokenUuid, nil
}

func (sys *TokenSystem) VerifyToken(token string) (*User, error) {
	var username string
	found := false

	err := sys.backend.db.View(func(txn *badger.Txn) error {
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
				return err
			}

			if found {
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	userSys := NewUserSystem(sys.backend)
	return userSys.GetUser(username)
}

func (sys *TokenSystem) ListTokens(username string) ([]string, error) {
	tokenUUIDs := []string{}

	err := sys.backend.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(fmt.Sprintf("pat:%s:", username))
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			tokenUUIDs = append(tokenUUIDs, string(key[len(username)+5:]))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tokenUUIDs, nil
}

func (sys *TokenSystem) DeleteToken(username string, tokenUUID string) error {
	return sys.backend.db.Update(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("pat:%s:%s", username, tokenUUID))
		err := txn.Delete(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return fmt.Errorf(
				"failed to delete token '%s' for user '%s': %w",
				tokenUUID, username, err,
			)
		}

		return nil
	})
}

func newToken(length int) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(alphabet)))

	for i := 0; i < length; i++ {
		r, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}

		result[i] = alphabet[r.Int64()]
	}

	return fmt.Sprintf("pat_%s", string(result)), nil
}
