package hash_test

import (
	"testing"

	"link-society.com/flowg/internal/utils/hash"
)

func TestHashPassword(t *testing.T) {
	h, err := hash.HashPassword("password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	ok, err := hash.VerifyPassword("password", h)
	if err != nil {
		t.Fatalf("VerifyPassword() error = %v", err)
	}

	if !ok {
		t.Fatal("VerifyPassword() = false; want true")
	}
}

func TestHashToken(t *testing.T) {
	h := hash.HashToken("pat_abc123")

	// SHA-256 hex is deterministic and 64 characters long.
	if got := hash.HashToken("pat_abc123"); got != h {
		t.Fatalf("HashToken() not deterministic: %q != %q", got, h)
	}
	if len(h) != 64 {
		t.Fatalf("HashToken() length = %d; want 64", len(h))
	}

	// Distinct tokens hash to distinct digests.
	if hash.HashToken("pat_other") == h {
		t.Fatal("HashToken() collided for distinct tokens")
	}
}

func BenchmarkHashPassword(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := hash.HashPassword("password")
		if err != nil {
			b.Fatalf("HashPassword() error = %v", err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	h, err := hash.HashPassword("password")
	if err != nil {
		b.Fatalf("HashPassword() error = %v", err)
	}

	for n := 0; n < b.N; n++ {
		ok, err := hash.VerifyPassword("password", h)
		if err != nil {
			b.Fatalf("VerifyPassword() error = %v", err)
		}

		if !ok {
			b.Fatal("VerifyPassword() = false; want true")
		}
	}
}
