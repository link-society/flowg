package vrl_test

import (
	"reflect"
	"testing"

	"link-society.com/flowg/internal/utils/langs/vrl"
)

func TestTransformLog(t *testing.T) {
	runner, err := vrl.NewScriptRunner(`
		.foo = "bar"
		.bar.baz = [1, 2, 3, "a"]
	`)
	if err != nil {
		t.Fatalf("NewScriptRunner() failed: %v", err)
	}
	defer runner.Close()

	input := map[string]string{}
	output, err := runner.TransformLog(input)
	if err != nil {
		t.Errorf("TransformLog() failed: %v", err)
	}

	expected := []map[string]string{{
		"foo":       "bar",
		"bar.baz.0": "1",
		"bar.baz.1": "2",
		"bar.baz.2": "3",
		"bar.baz.3": "a",
	}}
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("TransformLog() = %v, want %v", output, expected)
	}
}

func TestTransformLog_EmptyScript(t *testing.T) {
	runner, err := vrl.NewScriptRunner(``)
	if err != nil {
		t.Fatalf("NewScriptRunner() failed: %v", err)
	}
	defer runner.Close()

	input := map[string]string{
		"foo": "bar",
	}

	output, err := runner.TransformLog(input)
	if err != nil {
		t.Errorf("TransformLog() failed: %v", err)
	}

	if !reflect.DeepEqual(output, []map[string]string{input}) {
		t.Errorf("TransformLog() = %v, want %v", output, input)
	}
}
