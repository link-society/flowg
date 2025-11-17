type ForwarderConfigOtlpModel = {
  type: 'otlp'
  endpoint: string
  headers?: Record<string, string>
}

export default ForwarderConfigOtlpModel

export const factory = (): ForwarderConfigOtlpModel => ({
  type: 'otlp',
  endpoint: '',
  headers: undefined,
})
