type ForwarderConfigGoogleCloudLoggingModel = {
  type: 'googlecloudlogging'
  project_id: string
  log_id: string
  host: string
  port: number
  disable_tls: boolean
  disable_auth: boolean
  auth_json?: string
}

export default ForwarderConfigGoogleCloudLoggingModel

export const factory = (): ForwarderConfigGoogleCloudLoggingModel => ({
  type: 'googlecloudlogging',
  project_id: 'flowg',
  log_id: 'flowg',
  host: 'logging.googleapis.com',
  port: 443,
  disable_tls: false,
  disable_auth: false,
})
