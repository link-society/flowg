package hash_test

import (
	"testing"

	"link-society.com/flowg/internal/hash"
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
