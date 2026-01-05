package models

type ForwarderSyslogV2TagField string
type ForwarderSyslogV2SeverityField string
type ForwarderSyslogV2FacilityField string
type ForwarderSyslogV2MessageField string

func (ForwarderSyslogV2TagField) JSONSchemaOneOf() []any {
	return []any{
		string(""),
		DynamicField{},
	}
}

func (ForwarderSyslogV2SeverityField) JSONSchemaOneOf() []any {
	return []any{
		EnumField{"emerg", "alert", "crit", "err", "warning", "notice", "info", "debug"},
		DynamicField{},
	}
}

func (ForwarderSyslogV2FacilityField) JSONSchemaOneOf() []any {
	return []any{
		EnumField{"kern", "user", "mail", "daemon", "auth", "syslog", "lpr", "news", "uucp", "cron", "authpriv", "ftp", "local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7"},
		DynamicField{},
	}
}

func (ForwarderSyslogV2MessageField) JSONSchemaOneOf() []any {
	return []any{
		string(""),
		DynamicField{},
	}
}
