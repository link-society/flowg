type ForwarderConfigCloudWatchModel = {
  type: 'cloudwatch'
  app_id?: string
  endpoint: string
  region?: string
  access_key_id?: string
  secret_access_key?: string
  session_token?: string
  group: string
  stream: string
}

export default ForwarderConfigCloudWatchModel

export const factory = (): ForwarderConfigCloudWatchModel => ({
  type: 'cloudwatch',
  app_id: 'flowg',
  endpoint: '',
  region: '',
  access_key_id: '',
  secret_access_key: '',
  session_token: '',
  group: 'flowg',
  stream: 'logs',
})
