package vrl_test

import (
	"reflect"
	"testing"

	"link-society.com/flowg/internal/ffi/vrl"
)

func TestProcessRecord(t *testing.T) {
	input := map[string]string{}
	script := `
		.foo = "bar"
		.bar.baz = [1, 2, 3, "a"]
	`

	output, err := vrl.ProcessRecord(input, script)
	if err != nil {
		t.Errorf("ProcessRecord() failed: %v", err)
	}

	expected := map[string]string{
		"foo":       "bar",
		"bar.baz.0": "1",
		"bar.baz.1": "2",
		"bar.baz.2": "3",
		"bar.baz.3": "a",
	}
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("ProcessRecord() = %v, want %v", output, expected)
	}
}

func TestProcessRecord_EmptyScript(t *testing.T) {
	input := map[string]string{
		"foo": "bar",
	}
	script := ""

	output, err := vrl.ProcessRecord(input, script)
	if err != nil {
		t.Errorf("ProcessRecord() failed: %v", err)
	}

	if !reflect.DeepEqual(output, input) {
		t.Errorf("ProcessRecord() = %v, want %v", output, input)
	}
}
