import { DynamicField } from '@/lib/models/DynamicField.ts'

type ForwarderConfigAmqpModel = {
  type: 'amqp'
  url: string
  exchange: DynamicField<string>
  routing_key: DynamicField<string>
  body: DynamicField<string>
}

export default ForwarderConfigAmqpModel

export const factory = (): ForwarderConfigAmqpModel => ({
  type: 'amqp',
  url: '',
  exchange: '',
  routing_key: '',
  body: '@expr:toJSON(body)',
})
