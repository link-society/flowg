export type ClickhouseForwarderModel = {
	type: 'clickhouse'
	url: string
	db: string
	table: string
	user: string
	pass: string
	tls: boolean
}
