import LogEntryModel from '@/lib/models/LogEntryModel'

export type LogChartProps = Readonly<{
  rowData: LogEntryModel[]
  from: Date
  to: Date
}>
