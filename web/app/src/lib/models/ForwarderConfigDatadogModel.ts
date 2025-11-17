type ForwarderConfigDatadogModel = {
  type: 'datadog'
  url: string
  apiKey: string
}

export default ForwarderConfigDatadogModel

export const factory = (): ForwarderConfigDatadogModel => ({
  type: 'datadog',
  url: 'https://http-intake.logs.datadoghq.com/api/v2/logs',
  apiKey: '',
})
