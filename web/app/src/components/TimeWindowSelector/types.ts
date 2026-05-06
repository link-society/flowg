export type RelTimeWindowOption = {
  key: string
  label: string
  value: number
}

export type LabelRendererProps = Readonly<{
  timewindowType: 'relative' | 'absolute'
  relativeTimewindow: number
  from: Date
  to: Date
  live: boolean
}>

export type TimeWindow = {
  from: Date
  to: Date
  live: boolean
}

export type TimeWindowFactory = {
  make: () => TimeWindow
}

export type TimeWindowSelectorProps = Readonly<{
  onTimeWindowChanged: (factory: TimeWindowFactory) => void
}>
