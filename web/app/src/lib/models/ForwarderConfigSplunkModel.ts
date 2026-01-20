import { DynamicField } from '@/lib/models/DynamicField.ts'

type ForwarderConfigSplunkModel = {
  type: 'splunk'
  endpoint: string
  token: string
  source: DynamicField<string>
  host: DynamicField<string>
}

export default ForwarderConfigSplunkModel

export const factory = (): ForwarderConfigSplunkModel => ({
  type: 'splunk',
  endpoint: '',
  token: '',
  source: 'flowg',
  host: '@expr:log.host',
})
