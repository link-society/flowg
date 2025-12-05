type PermissionsModel = {
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

  can_read_system_configuration: boolean
  can_write_system_configuration: boolean
}

export default PermissionsModel
