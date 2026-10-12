package auth

import (
	"fmt"
	"strings"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

// DefaultTokenTTL is how long a JWT issued on login stays valid.
const DefaultTokenTTL = 24 * time.Hour

// Keys holds the secret material used to sign and verify the JWTs issued on
// login. It is provided to the dependency graph from the configuration.
//
// Every node of a cluster must share the same SigningKey, otherwise a token
// issued by one node is rejected by the others.
type Keys struct {
	// SigningKey is the HMAC secret signing and verifying the JWTs.
	SigningKey []byte
	// TokenTTL is how long an issued JWT stays valid.
	TokenTTL time.Duration
}

// NewJWT issues a signed, time-limited token that proves the bearer
// authenticated as the given user.
//
// It is the credential handed back on successful login, letting subsequent
// requests authenticate without replaying the password. The returned string
// carries a "jwt_" prefix so that [ApiMiddleware] can tell it apart from other
// bearer credentials, and the underlying token expires after [Keys.TokenTTL]
// so that a leaked credential is only useful for a bounded window.
func (k Keys) NewJWT(username string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "flowg",
		"sub": username,
		"exp": time.Now().Add(k.TokenTTL).Unix(),
	})

	tok, err := t.SignedString(k.SigningKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return fmt.Sprintf("jwt_%s", tok), nil
}

// VerifyJWT authenticates a token previously issued by [Keys.NewJWT] and
// recovers the user it was minted for.
//
// It is the counterpart consumed by [ApiMiddleware] to turn a bearer credential
// back into an identity. It accepts a token only when it carries the expected
// "jwt_" prefix, its signature matches [Keys.SigningKey], and it has not
// expired; any other case is reported as an error rather than a best-effort
// identity, so an invalid credential can never be mistaken for a valid one.
func (k Keys) VerifyJWT(token string) (string, error) {
	if !strings.HasPrefix(token, "jwt_") {
		return "", fmt.Errorf("invalid token prefix")
	}
	token = token[len("jwt_"):]

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return k.SigningKey, nil
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
