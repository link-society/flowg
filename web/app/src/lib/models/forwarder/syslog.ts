export type SyslogForwarderModel = {
  type: 'syslog'
  network: 'tcp' | 'udp'
  address: string
  tag: string
  severity: SyslogSeverity
  facility: SyslogFacility
}

export const SyslogSeverityValues = [
  'emerg',
  'alert',
  'crit',
  'err',
  'warning',
  'notice',
  'info',
  'debug',
]

export type SyslogSeverity = typeof SyslogSeverityValues[number]

export const SyslogFacilityValues = [
  'kern',
  'user',
  'mail',
  'daemon',
  'auth',
  'syslog',
  'lpr',
  'news',
  'uucp',
  'cron',
  'authpriv',
  'ftp',
  'local0',
  'local1',
  'local2',
  'local3',
  'local4',
  'local5',
  'local6',
  'local7',
]

export type SyslogFacility = typeof SyslogFacilityValues[number]
