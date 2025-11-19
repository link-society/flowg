export type ClickhouseForwarderModel = {
  type: 'clickhouse'
  address: string
  db: string
  table: string
  user: string
  pass: string
  tls: boolean
}
