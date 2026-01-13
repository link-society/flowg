import { DynamicField } from '@/lib/models/DynamicField'

type ForwarderConfigDatadogModel = {
  type: 'datadog'
  url: string
  apiKey: string
  ddsource: DynamicField<string>
  ddtags: DynamicField<string>
  hostname: DynamicField<string>
  message: DynamicField<string>
  service: DynamicField<string>
}

export default ForwarderConfigDatadogModel

export const factory = (): ForwarderConfigDatadogModel => ({
  type: 'datadog',
  url: 'https://http-intake.logs.datadoghq.com/api/v2/logs',
  apiKey: '',
  ddsource: '@expr:log.ddsource',
  ddtags: '@expr:log.ddtags',
  hostname: '@expr:log.hostname',
  message: '@expr:log.message',
  service: '@expr:log.service',
})
