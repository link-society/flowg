package models

import "github.com/swaggest/jsonschema-go"

type ForwarderSyslogV2TagField string
type ForwarderSyslogV2SeverityField string
type ForwarderSyslogV2FacilityField string
type ForwarderSyslogV2MessageField string

type ForwarderSyslogV2SeverityEnumType string
type ForwarderSyslogV2FacilityEnumType string

func (ForwarderSyslogV2TagField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderSyslogV2SeverityField) JSONSchemaAnyOf() []any {
	return []any{
		ForwarderSyslogV2SeverityEnumType(""),
		DynamicField(""),
	}
}

func (ForwarderSyslogV2FacilityField) JSONSchemaAnyOf() []any {
	return []any{
		ForwarderSyslogV2FacilityEnumType(""),
		DynamicField(""),
	}
}

func (ForwarderSyslogV2MessageField) JSONSchemaAnyOf() []any {
	return []any{
		string(""),
		DynamicField(""),
	}
}

func (ForwarderSyslogV2SeverityEnumType) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithEnum("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug")
	return nil
}

func (ForwarderSyslogV2FacilityEnumType) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.WithEnum("kern", "user", "mail", "daemon", "auth", "syslog", "lpr", "news", "uucp", "cron", "authpriv", "ftp", "local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7")
	return nil
}
