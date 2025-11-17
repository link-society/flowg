type ForwarderConfigElasticModel = {
  type: 'elastic'
  index: string
  username: string
  password: string
  addresses: string[]
  ca?: string
}

export default ForwarderConfigElasticModel

export const factory = (): ForwarderConfigElasticModel => ({
  type: 'elastic',
  index: 'default',
  username: '',
  password: '',
  addresses: ['https://localhost:9200'],
  ca: undefined,
})
