type ForwarderConfigSplunkModel = {
  type: 'splunk'
  endpoint: string
  token: string
}

export default ForwarderConfigSplunkModel

export const factory = (): ForwarderConfigSplunkModel => ({
  type: 'splunk',
  endpoint: '',
  token: '',
})
