package auth_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"link-society.com/flowg/api/auth"
)

func TestNewJWTRoundTrip(t *testing.T) {
	// Contract: a token issued by NewJWT verifies back to the same subject.
	token, err := auth.NewJWT("alice")
	require.NoError(t, err)

	username, err := auth.VerifyJWT(token)
	require.NoError(t, err)
	assert.Equal(t, "alice", username)
}

func TestNewJWTCarriesPrefix(t *testing.T) {
	// Contract: issued tokens are prefixed so the middleware can tell them
	// apart from other bearer credentials.
	token, err := auth.NewJWT("alice")
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(token, "jwt_"), "token must carry the jwt_ prefix")
}

func TestVerifyJWTRejectsMissingPrefix(t *testing.T) {
	// Contract: a token without the expected prefix is not a JWT and must be
	// rejected.
	_, err := auth.VerifyJWT("alice")
	require.Error(t, err)
}

func TestVerifyJWTRejectsTamperedToken(t *testing.T) {
	// Contract: a token whose signature does not match the signing key is
	// rejected rather than yielding a best-effort identity.
	token, err := auth.NewJWT("alice")
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

	_, err = auth.VerifyJWT(tampered)
	require.Error(t, err)
}

func TestVerifyJWTRejectsForeignSignature(t *testing.T) {
	// Contract: a well-formed JWT signed with a different key is rejected.
	// "jwt_<header>.<payload>.<sig>" signed with an unrelated secret.
	foreign := "jwt_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
		"eyJpc3MiOiJmbG93ZyIsInN1YiI6ImFsaWNlIn0." +
		"3Vd3X8mO0p1q2r3s4t5u6v7w8x9y0zAbCdEfGhIjKlM"

	_, err := auth.VerifyJWT(foreign)
	require.Error(t, err)
}
