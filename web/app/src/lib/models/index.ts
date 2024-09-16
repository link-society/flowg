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

  can_view_alerts: boolean
  can_edit_alerts: boolean

  can_view_acls: boolean
  can_edit_acls: boolean

  can_send_logs: boolean
}

export type StreamConfigModel = {
  indexed_fields: string[]
  ttl: number
  size: number
}

export type TokenModel = {
  token: string
  token_uuid: string
}
