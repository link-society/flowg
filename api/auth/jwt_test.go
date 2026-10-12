package auth_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"link-society.com/flowg/api/auth"
)

var testKeys = auth.Keys{
	SigningKey: []byte("test-signing-key"),
	TokenTTL:   auth.DefaultTokenTTL,
}

func TestNewJWTRoundTrip(t *testing.T) {
	// Contract: a token issued by Keys.NewJWT verifies back to the same subject.
	token, err := testKeys.NewJWT("alice")
	require.NoError(t, err)

	username, err := testKeys.VerifyJWT(token)
	require.NoError(t, err)
	assert.Equal(t, "alice", username)
}

func TestNewJWTCarriesPrefix(t *testing.T) {
	// Contract: issued tokens are prefixed so the middleware can tell them
	// apart from other bearer credentials.
	token, err := testKeys.NewJWT("alice")
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(token, "jwt_"), "token must carry the jwt_ prefix")
}

func TestVerifyJWTRejectsMissingPrefix(t *testing.T) {
	// Contract: a token without the expected prefix is not a JWT and must be
	// rejected.
	_, err := testKeys.VerifyJWT("alice")
	require.Error(t, err)
}

func TestVerifyJWTRejectsTamperedToken(t *testing.T) {
	// Contract: a token whose signature does not match the signing key is
	// rejected rather than yielding a best-effort identity.
	token, err := testKeys.NewJWT("alice")
	require.NoError(t, err)

	// Break the signature by mutating the first character of its segment.
	// The final base64url character of a 32-byte HMAC signature only carries
	// 4 significant bits (the trailing 2 are discarded on decode), so flipping
	// it can decode to the same bytes and leave the signature valid. The first
	// character's bits are always significant, guaranteeing a real change.
	parts := strings.Split(strings.TrimPrefix(token, "jwt_"), ".")
	require.Len(t, parts, 3)

	sig := parts[2]
	if sig[0] == 'A' {
		parts[2] = "B" + sig[1:]
	} else {
		parts[2] = "A" + sig[1:]
	}
	tampered := "jwt_" + strings.Join(parts, ".")

	_, err = testKeys.VerifyJWT(tampered)
	require.Error(t, err)
}

func TestVerifyJWTRejectsForeignSignature(t *testing.T) {
	// Contract: a well-formed JWT signed with a different key is rejected.
	// "jwt_<header>.<payload>.<sig>" signed with an unrelated secret.
	foreign := "jwt_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
		"eyJpc3MiOiJmbG93ZyIsInN1YiI6ImFsaWNlIn0." +
		"3Vd3X8mO0p1q2r3s4t5u6v7w8x9y0zAbCdEfGhIjKlM"

	_, err := testKeys.VerifyJWT(foreign)
	require.Error(t, err)
}

func TestVerifyJWTRejectsTokenSignedWithOtherKeys(t *testing.T) {
	// Contract: a token issued with other keys (e.g. by a node configured with
	// another secret key) is rejected.
	otherKeys := auth.Keys{
		SigningKey: []byte("other-signing-key"),
		TokenTTL:   auth.DefaultTokenTTL,
	}

	token, err := otherKeys.NewJWT("alice")
	require.NoError(t, err)

	_, err = testKeys.VerifyJWT(token)
	require.Error(t, err)
}

func TestVerifyJWTRejectsExpiredToken(t *testing.T) {
	// Contract: a token is only valid for the keys' token TTL.
	expiringKeys := auth.Keys{
		SigningKey: testKeys.SigningKey,
		TokenTTL:   -time.Minute,
	}

	token, err := expiringKeys.NewJWT("alice")
	require.NoError(t, err)

	_, err = testKeys.VerifyJWT(token)
	require.Error(t, err)
}
