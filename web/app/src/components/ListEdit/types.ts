type GetKey = (item: string, index: number) => string

export type ListEditProps = Readonly<{
  id: string
  list: Array<string>
  setList: (list: Array<string>) => void
  getKey?: GetKey
}>
