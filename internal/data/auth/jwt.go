package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SIGNING_KEY []byte

func init() {
	key, err := newToken(32)
	if err != nil {
		panic(err)
	}

	JWT_SIGNING_KEY = []byte(key)
}

func NewJWT(username string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "flowg",
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return t.SignedString([]byte(JWT_SIGNING_KEY))
}

func VerifyJWT(token string) (string, error) {
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
