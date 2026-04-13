export type ValueChipProps = {
  value: string
  selected: boolean
  onToggle: (selected: boolean) => void
}

export type StreamIndexSelectorProps = {
  indices: Record<string, Array<string>>
  selection: Record<string, Array<string>>
  onSelectionChange: (selection: Record<string, Array<string>>) => void
}
