type ForwarderConfigHttpModel = {
  type: 'http'
  url: string
  headers?: Record<string, string>
}

export default ForwarderConfigHttpModel

export const factory = (): ForwarderConfigHttpModel => ({
  type: 'http',
  url: '',
  headers: undefined,
})
