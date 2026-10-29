package models

import "github.com/swaggest/jsonschema-go"

// ForwarderSyslogV2 forwards records to a syslog server over TCP or UDP. Tag,
// severity, facility and message are dynamic fields evaluated per record.
type ForwarderSyslogV2 struct {
	Type     string                         `json:"type" enum:"syslog" required:"true"`
	Network  string                         `json:"network" enum:"tcp,udp" required:"true"`
	Address  string                         `json:"address" required:"true" pattern:"^(([a-zA-Z0-9.-]+)|(\\[[0-9A-Fa-f:]+\\])):[0-9]{1,5}$"`
	Tag      ForwarderSyslogV2TagField      `json:"tag" required:"true"`
	Severity ForwarderSyslogV2SeverityField `json:"severity" required:"true"`
	Facility ForwarderSyslogV2FacilityField `json:"facility" required:"true"`
	Message  ForwarderSyslogV2MessageField  `json:"message" required:"false"`
}

// The Forwarder*Field types below are the syslog forwarder's per-record fields:
// each is either a literal value or a DynamicField ("@expr:" expression), as
// advertised by their JSONSchemaAnyOf methods. The *EnumType helpers constrain
// the literal form of severity and facility to their allowed keywords.
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
