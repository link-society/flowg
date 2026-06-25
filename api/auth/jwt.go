package auth

import (
	"fmt"
	"strings"

	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"link-society.com/flowg/internal/storage/backends/badger/concrete/auth/secret"
)

// JWT_SIGNING_KEY is the secret used to sign and verify the JWTs minted by this
// package.
//
// It is sourced once, at startup, from the FLOWG_SECRET_KEY environment
// variable. When the variable is unset a random key is generated, which is
// convenient for local development but makes issued tokens valid only for the
// lifetime of the process — set FLOWG_SECRET_KEY explicitly in production so
// tokens survive restarts.
var JWT_SIGNING_KEY []byte

func init() {
	key := os.Getenv("FLOWG_SECRET_KEY")
	if key == "" {
		var err error
		key, err = secret.NewSecret("jwt", 32)
		if err != nil {
			panic(err)
		}
	}

	JWT_SIGNING_KEY = []byte(key)
}

// NewJWT issues a signed, time-limited token that proves the bearer
// authenticated as the given user.
//
// It is the credential handed back on successful login, letting subsequent
// requests authenticate without replaying the password. The returned string
// carries a "jwt_" prefix so that [ApiMiddleware] can tell it apart from other
// bearer credentials, and the underlying token expires so that a leaked
// credential is only useful for a bounded window.
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

// VerifyJWT authenticates a token previously issued by [NewJWT] and recovers the
// user it was minted for.
//
// It is the counterpart consumed by [ApiMiddleware] to turn a bearer credential
// back into an identity. It accepts a token only when it carries the expected
// "jwt_" prefix, its signature matches [JWT_SIGNING_KEY], and it has not
// expired; any other case is reported as an error rather than a best-effort
// identity, so an invalid credential can never be mistaken for a valid one.
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
