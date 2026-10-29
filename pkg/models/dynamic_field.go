package models

import "github.com/swaggest/jsonschema-go"

// DynamicField is a forwarder configuration value that may either be a literal
// string or an [expr](https://expr-lang.org/) expression (when prefixed with
// "@expr:"), evaluated against the record at forward time.
type DynamicField string

// PrepareJSONSchema describes a DynamicField for OpenAPI generation: a string
// matching the "@expr:" convention.
func (DynamicField) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.AddType(jsonschema.String)
	schema.WithPattern(`^@expr:`)
	return nil
}
