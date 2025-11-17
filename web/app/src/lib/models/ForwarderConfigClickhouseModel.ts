type ForwarderConfigClickhouseModel = {
  type: 'clickhouse'
  address: string
  db: string
  table: string
  user: string
  pass: string
  tls: boolean
}

export default ForwarderConfigClickhouseModel

export const factory = (): ForwarderConfigClickhouseModel => ({
  type: 'clickhouse',
  address: 'localhost:9000',
  db: 'default',
  table: 'default',
  user: 'default',
  pass: '',
  tls: true,
})
