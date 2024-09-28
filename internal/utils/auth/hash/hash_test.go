package hash_test

import (
	"testing"

	"link-society.com/flowg/internal/utils/auth/hash"
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
