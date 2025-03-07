import { useState } from 'react'

import Paper from '@mui/material/Paper'
import Drawer from '@mui/material/Drawer'

import { AgGridReact } from 'ag-grid-react'
import { ColDef } from 'ag-grid-community'

import { LogEntryModel } from '@/lib/models'

type LogTableProps = Readonly<{
  rowData: LogEntryModel[]
  columnDefs: ColDef<LogEntryModel>[]
}>

export const LogTable = (props: LogTableProps) => {
  const [selectedRow, setSelectedRow] = useState<LogEntryModel | undefined>(undefined)

  return (
    <Paper className="grow ag-theme-balham">
      <AgGridReact<LogEntryModel>
        rowData={props.rowData}
        columnDefs={props.columnDefs}
        suppressFieldDotNotation
        enableCellTextSelection
        autoSizeStrategy={{type: 'fitCellContents'}}
        onRowDoubleClicked={(e) => {
          setSelectedRow(e.data)
        }}
      />
      <Drawer
        anchor="right"
        open={selectedRow !== undefined}
        onClose={() => setSelectedRow(undefined)}
        sx={{
          '& .MuiDrawer-paper': {
            width: '33vw',
            padding: '0.75rem',
          },
        }}
      >
        <Paper
          variant="outlined"
          className="
            p-2 w-full overflow-auto
            font-mono !bg-gray-100
          "
          component="pre"
        >
          {JSON.stringify(selectedRow, null, 2)}
        </Paper>
      </Drawer>
    </Paper>
  )
}
