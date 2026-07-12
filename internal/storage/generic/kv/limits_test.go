package kv_test

import (
	"errors"
	"testing"

	"link-society.com/flowg/internal/storage/generic/kv"
)

func TestCheckKeySize(t *testing.T) {
	if err := kv.CheckKeySize(kv.MaxKeySize); err != nil {
		t.Fatalf("expected key of exactly MaxKeySize to be accepted, got %v", err)
	}

	err := kv.CheckKeySize(kv.MaxKeySize + 1)
	if !errors.Is(err, kv.ErrKeyTooLarge) {
		t.Fatalf("expected ErrKeyTooLarge, got %v", err)
	}
}

func TestCheckValueSize(t *testing.T) {
	if err := kv.CheckValueSize(kv.MaxValueSize); err != nil {
		t.Fatalf("expected value of exactly MaxValueSize to be accepted, got %v", err)
	}

	err := kv.CheckValueSize(kv.MaxValueSize + 1)
	if !errors.Is(err, kv.ErrValueTooLarge) {
		t.Fatalf("expected ErrValueTooLarge, got %v", err)
	}
}
