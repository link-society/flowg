import { ColDef } from 'ag-grid-community'

import LogEntryModel from '@/lib/models/LogEntryModel'

export type LogTableProps = Readonly<{
  rowData: LogEntryModel[]
  columnDefs: ColDef<LogEntryModel>[]
}>
