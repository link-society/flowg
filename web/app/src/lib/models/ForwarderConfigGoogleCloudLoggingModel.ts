type ForwarderConfigGoogleCloudLoggingModel = {
  type: 'googlecloudlogging'
  project_id: string
  log_id: string
  endpoint: string
  auth_json?: string
}

export default ForwarderConfigGoogleCloudLoggingModel

export const factory = (): ForwarderConfigGoogleCloudLoggingModel => ({
  type: 'googlecloudlogging',
  project_id: 'flowg',
  log_id: 'flowg',
  endpoint: 'logging.googleapis.com',
})
