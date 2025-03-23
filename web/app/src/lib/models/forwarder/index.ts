import { HttpForwarderModel } from '@/lib/models/forwarder/http'

export type ForwarderModel = {
  config: ForwarderConfigModel
}

export const ForwarderTypeValues = [
  { key: 'http', label: 'Webhook' },
] as const

export type ForwarderConfigModel =
  | HttpForwarderModel

export type ForwarderTypes = ForwarderConfigModel['type']
