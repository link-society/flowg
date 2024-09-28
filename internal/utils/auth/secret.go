package auth

import (
	"fmt"

	"crypto/rand"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewSecret(prefix string, length int) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(alphabet)))

	for i := 0; i < length; i++ {
		r, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}

		result[i] = alphabet[r.Int64()]
	}

	return fmt.Sprintf("%s_%s", prefix, string(result)), nil
}
