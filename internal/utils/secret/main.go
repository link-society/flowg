package secret

import (
	"fmt"

	"crypto/rand"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var alphabetSize = big.NewInt(int64(len(alphabet)))

// NewSecret generates a new secret with the given prefix and length.
// The secret is composed of random characters from a predefined alphabet.
func NewSecret(prefix string, length int) (string, error) {
	result := make([]byte, length)

	for i := range result {
		r, err := rand.Int(rand.Reader, alphabetSize)
		if err != nil {
			return "", err
		}

		result[i] = alphabet[r.Int64()]
	}

	return fmt.Sprintf("%s_%s", prefix, string(result)), nil
}
