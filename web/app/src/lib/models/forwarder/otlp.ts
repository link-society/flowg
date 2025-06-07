export type OtlpForwarderModel = {
  type: 'otlp'
  endpoint: string
  headers?: Record<string, string>
}
