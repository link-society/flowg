import { ReactFlowJsonObject } from '@xyflow/react'

export type ProfileModel = {
  user: UserModel
  permissions: PermissionsModel
}

export type UserModel = {
  name: string
  roles: string[]
}

export type PermissionsModel = {
  can_view_pipelines: boolean
  can_edit_pipelines: boolean

  can_view_transformers: boolean
  can_edit_transformers: boolean

  can_view_streams: boolean
  can_edit_streams: boolean

  can_view_forwarders: boolean
  can_edit_forwarders: boolean

  can_view_acls: boolean
  can_edit_acls: boolean

  can_send_logs: boolean
}

export type RoleModel = {
  name: string
  scopes: string[]
}

export type StreamConfigModel = {
  indexed_fields: string[]
  ttl: number
  size: number
}

export type ForwarderModel = {
  config: ForwarderConfigModel
}

export const ForwarderTypeValues = [
  { key: 'http', label: 'Webhook' },
] as const

export type ForwarderConfigModel =
  | ForwarderConfigHttpModel

export type ForwarderConfigHttpModel = {
  type: 'http'
  url: string
  headers: Record<string, string>
}

export type ForwarderTypes = ForwarderConfigModel['type']

export type TokenModel = {
  token: string
  token_uuid: string
}

export type PipelineModel = {
  nodes: ReactFlowJsonObject['nodes'],
  edges: ReactFlowJsonObject['edges'],
}

export type LogEntryModel = {
  timestamp: Date
  fields: Record<string, string>
}
