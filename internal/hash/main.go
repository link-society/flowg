package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"strings"

	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type hashParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func HashPassword(password string) (string, error) {
	hp := hashParams{
		memory:      1, // 1 kB
		iterations:  1,
		parallelism: 1,
		saltLength:  4,
		keyLength:   32,
	}

	salt := make([]byte, hp.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", &HashError{Reason: fmt.Errorf("failed to generate hash: %w", err)}
	}

	key := argon2.IDKey(
		[]byte(password),
		salt,
		hp.iterations,
		hp.memory,
		hp.parallelism,
		hp.keyLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		hp.memory,
		hp.iterations,
		hp.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func VerifyPassword(password, hash string) (bool, error) {
	hp, salt, key, err := decodeHash(hash)
	if err != nil {
		return false, err
	}

	candidateKey := argon2.IDKey(
		[]byte(password),
		salt,
		hp.iterations,
		hp.memory,
		hp.parallelism,
		hp.keyLength,
	)

	return subtle.ConstantTimeCompare(key, candidateKey) == 1, nil
}

func decodeHash(hash string) (*hashParams, []byte, []byte, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, &HashError{Reason: fmt.Errorf("invalid hash format")}
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, &HashError{Reason: fmt.Errorf("failed to decode hash: %w", err)}
	}
	if version != argon2.Version {
		return nil, nil, nil, &HashError{Reason: fmt.Errorf("incompatible version")}
	}

	hp := &hashParams{}
	_, err = fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&hp.memory,
		&hp.iterations,
		&hp.parallelism,
	)
	if err != nil {
		return nil, nil, nil, &HashError{Reason: fmt.Errorf("failed to decode hash: %w", err)}
	}

	salt, err := base64.RawStdEncoding.DecodeString(string(parts[4]))
	if err != nil {
		return nil, nil, nil, &HashError{Reason: fmt.Errorf("failed to decode salt: %w", err)}
	}
	hp.saltLength = uint32(len(salt))

	key, err := base64.RawStdEncoding.DecodeString(string(parts[5]))
	if err != nil {
		return nil, nil, nil, &HashError{Reason: fmt.Errorf("failed to decode key: %w", err)}
	}
	hp.keyLength = uint32(len(key))

	return hp, []byte(salt), []byte(key), nil
}
