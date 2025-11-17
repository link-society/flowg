type ForwarderConfigAmqpModel = {
  type: 'amqp'
  url: string
  exchange: string
  routing_key: string
}

export default ForwarderConfigAmqpModel

export const factory = (): ForwarderConfigAmqpModel => ({
  type: 'amqp',
  url: '',
  exchange: '',
  routing_key: '',
})
