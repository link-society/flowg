import { DynamicField } from '@/lib/models/DynamicField'

type ForwarderConfigSyslogModel = {
  type: 'syslog'
  network: SyslogNetwork
  address: string
  tag: DynamicField<string>
  severity: DynamicField<SyslogSeverity>
  facility: DynamicField<SyslogFacility>
  message: DynamicField<string>
}

export const SyslogNetworkValues = ['tcp', 'udp'] as const

export type SyslogNetwork = (typeof SyslogNetworkValues)[number]

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

export type SyslogSeverity = (typeof SyslogSeverityValues)[number]

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

export type SyslogFacility = (typeof SyslogFacilityValues)[number]

export default ForwarderConfigSyslogModel

export const factory = (): ForwarderConfigSyslogModel => ({
  type: 'syslog',
  network: 'tcp',
  address: '',
  tag: '',
  severity: 'info',
  facility: 'user',
  message: '@expr:toJSON(log)',
})
