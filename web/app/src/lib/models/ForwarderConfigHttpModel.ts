import { DynamicField } from '@/lib/models/DynamicField.ts'

type ForwarderConfigHttpModel = {
  type: 'http'
  url: string
  headers?: Record<string, string>
  body: DynamicField<string>
}

export default ForwarderConfigHttpModel

export const factory = (): ForwarderConfigHttpModel => ({
  type: 'http',
  url: '',
  headers: undefined,
  body: '@expr:toJSON(log)',
})
