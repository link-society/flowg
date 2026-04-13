import { ReactNode } from 'react'

export type InputTransferListProps<T> = Readonly<{
  choices: T[]
  getItemId: (item: T) => string
  renderItem: (item: T) => ReactNode
  onChoiceUpdate: (choices: readonly T[]) => void
}>
