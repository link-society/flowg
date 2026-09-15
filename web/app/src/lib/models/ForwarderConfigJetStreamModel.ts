type ForwarderConfigJetStreamModel = {
  type: 'jetstream'
  servers: string[]
  subject: string
  body: string
  message_id: string
  expected_stream: string
  publish_timeout: number
  headers?: Record<string, string>
  auth_mode: string
  username?: string
  password?: string
  token?: string
  credentials?: string
  tls_ca?: string
  tls_certificate?: string
  tls_private_key?: string
  tls_server_name?: string
  tls_insecure_skip_verify?: boolean
}

export type JetStreamAuth = (typeof JetStreamAuthValues)[number]

export const JetStreamAuthValues = [
  'none',
  'user_password',
  'token',
  'credentials',
]

export default ForwarderConfigJetStreamModel

export const factory = (): ForwarderConfigJetStreamModel => ({
  type: 'jetstream',
  servers: [],
  subject: '',
  body: '@expr:toJson(log)',
  message_id: '',
  expected_stream: '',
  publish_timeout: 0,
  auth_mode: JetStreamAuthValues[0],
})
