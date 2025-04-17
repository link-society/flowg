package auth

import (
	"fmt"
	"strings"

	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SIGNING_KEY []byte

func init() {
	key := os.Getenv("FLOWG_SECRET_KEY")
	if key == "" {
		var err error
		key, err = NewSecret("jwt", 32)
		if err != nil {
			panic(err)
		}
	}

	JWT_SIGNING_KEY = []byte(key)
}

func NewJWT(username string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "flowg",
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tok, err := t.SignedString([]byte(JWT_SIGNING_KEY))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return fmt.Sprintf("jwt_%s", tok), nil
}

func VerifyJWT(token string) (string, error) {
	if !strings.HasPrefix(token, "jwt_") {
		return "", fmt.Errorf("invalid token prefix")
	}
	token = token[len("jwt_"):]

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWT_SIGNING_KEY), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse JWT: %w", err)
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return "", fmt.Errorf("invalid JWT")
	}

	return claims["sub"].(string), nil
}
