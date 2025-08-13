export type ElasticForwarderModel = {
  type: 'elastic'
  index: string
  username: string
  password: string
  addresses: string[]
  ca?: string
}
