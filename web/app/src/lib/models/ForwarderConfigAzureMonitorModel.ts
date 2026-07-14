type ForwarderConfigAzureMonitorModel = {
  type: 'azuremonitor'
  endpoint: string
  token: string
  expires_on: string
  rule_id: string
  stream_name: string
  allow_insecure: boolean
}

export default ForwarderConfigAzureMonitorModel

export const factory = (): ForwarderConfigAzureMonitorModel => ({
  type: 'azuremonitor',
  endpoint: '',
  token: '',
  expires_on: '',
  rule_id: '',
  stream_name: '',
  allow_insecure: false,
})
