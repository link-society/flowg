package models_test

import (
	"testing"

	"encoding/json"
	"reflect"

	"github.com/expr-lang/expr"

	"github.com/swaggest/jsonschema-go"

	"link-society.com/flowg/internal/models"
)

func TestDynamicField_JSONSchemaShape(t *testing.T) {
	dynField := models.DynamicField("")

	var s jsonschema.Schema
	if err := dynField.PrepareJSONSchema(&s); err != nil {
		t.Fatalf("failed to prepare schema: %v", err)
	}

	raw, err := json.Marshal(&s)
	if err != nil {
		t.Fatalf("failed to marshal schema: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("failed to unmarshal schema: %v", err)
	}

	expected := map[string]any{
		"type":    "string",
		"pattern": "^@expr:",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("unexpected schema shape.\nGot:      %+v\nExpected: %+v", got, expected)
	}
}

func TestCompileDynamicField_Static(t *testing.T) {
	prog, err := models.CompileDynamicField("foobar")
	if err != nil {
		t.Fatalf("failed to compile dynamic field: %v", err)
	}

	output, err := expr.Run(prog, nil)
	if err != nil {
		t.Fatalf("failed to run compiled program: %v", err)
	}

	expected := "foobar"
	if output != expected {
		t.Fatalf("unexpected output.\nGot:      %q\nExpected: %q", output, expected)
	}
}

func TestCompileDynamicField_Expression(t *testing.T) {
	prog, err := models.CompileDynamicField("@expr:foo")
	if err != nil {
		t.Fatalf("failed to compile dynamic field: %v", err)
	}

	env := map[string]any{
		"foo": "bar",
	}
	output, err := expr.Run(prog, env)
	if err != nil {
		t.Fatalf("failed to run compiled program: %v", err)
	}

	expected := "bar"
	if output != expected {
		t.Fatalf("unexpected output.\nGot:      %q\nExpected: %q", output, expected)
	}
}
