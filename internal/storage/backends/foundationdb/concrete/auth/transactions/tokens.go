package transactions

import (
	"fmt"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/google/uuid"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/backends/foundationdb/concrete/auth/hash"
	"link-society.com/flowg/internal/storage/backends/foundationdb/concrete/auth/secret"
)

// Subspace layout for tokens:
//
//	<root>/pat/<username>/<uuid> → token hash  — PAT storage
var patSub subspace.Subspace

func initTokenSubspaces(root subspace.Subspace) {
	patSub = root.Sub("pat")
}

// CreateToken issues a fresh PAT for an existing user. It confirms the user
// exists through their index marker, then stores the token hash under a new
// <root>/pat/<username>/<uuid> key and returns the plaintext token (shown to
// the caller this one time only) along with its UUID.
func CreateToken(tr fdb.Transaction, username string) (string, string, error) {
	token, err := secret.NewSecret("pat", 32)
	if err != nil {
		return "", "", err
	}

	tokenHash, err := hash.HashPassword(token)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash token: %w", err)
	}

	tokenUuid := uuid.New().String()

	// Verify user exists
	_, err = tr.Get(indexUserSub.Pack(tuple.Tuple{username})).Get()
	if err != nil {
		return "", "", fmt.Errorf("user '%s' does not exist", username)
	}

	tr.Set(patSub.Sub(username).Pack(tuple.Tuple{tokenUuid}), []byte(tokenHash))

	return token, tokenUuid, nil
}

// VerifyToken resolves a plaintext token back to its owning user. Tokens are
// not addressable by their value, so it scans every <root>/pat/ key and
// argon2id-compares the candidate against each stored hash; on a match it
// returns the owning user.
func VerifyToken(tr fdb.ReadTransaction, token string) (*models.User, error) {
	var username string
	found := false

	iter := tr.GetRange(patSub, fdb.RangeOptions{}).Iterator()
	for iter.Advance() {
		kv := iter.MustGet()

		t, err := patSub.Unpack(kv.Key)
		if err != nil {
			continue
		}

		associatedUser := t[0].(string)
		tokenHash := string(kv.Value)

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

	return FetchUser(tr, username)
}

// ListTokens returns the UUIDs of every PAT owned by a user by scanning the
// <root>/pat/<username>/ prefix.
func ListTokens(tr fdb.ReadTransaction, username string) []string {
	tokenUUIDs := []string{}

	sub := patSub.Sub(username)
	iter := tr.GetRange(sub, fdb.RangeOptions{}).Iterator()
	for iter.Advance() {
		kv := iter.MustGet()

		t, err := sub.Unpack(kv.Key)
		if err != nil {
			continue
		}

		tokenUUIDs = append(tokenUUIDs, t[0].(string))
	}

	return tokenUUIDs
}

// DeleteToken removes a single PAT, identified by its UUID, from a user.
func DeleteToken(tr fdb.Transaction, username string, tokenUUID string) error {
	tr.Clear(patSub.Sub(username).Pack(tuple.Tuple{tokenUUID}))
	return nil
}
