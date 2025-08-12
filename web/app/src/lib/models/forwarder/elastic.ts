export type ElasticForwarderModel = {
  type: 'elastic'
  index: string
  addresses: string[]
  ca?: string
  token?: string
}
