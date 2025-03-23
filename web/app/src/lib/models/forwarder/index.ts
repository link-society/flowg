import { HttpForwarderModel } from '@/lib/models/forwarder/http'
import { SyslogForwarderModel } from '@/lib/models/forwarder/syslog'

export type ForwarderModel = {
  config: ForwarderConfigModel
}

export const ForwarderTypeValues = [
  { key: 'http', label: 'Webhook' },
  { key: 'syslog', label: 'Syslog Server' },
] as const

export const ForwarderTypeLabelMap = ForwarderTypeValues.reduce(
  (acc, cur) => {
    acc[cur.key] = cur.label
    return acc
  },
  {} as Record<ForwarderTypes, string>,
)

export type ForwarderConfigModel =
  | HttpForwarderModel
  | SyslogForwarderModel

export type ForwarderTypes = ForwarderConfigModel['type']
