import * as validators from '@/lib/validators'

export type InputListProps = {
  id?: string
  itemLabel?: string
  items: string[]
  itemValidators?: Array<validators.Validator<string>>
  onChange: (items: string[]) => void
}

export type Row = {
  id: string
  value: string
}
