import { useState } from 'react'

import { AgGridReact } from 'ag-grid-react'

import LogEntryModel from '@/lib/models/LogEntryModel'

import { LogTableContainer, LogTableDetailPre, LogTableDrawer } from './styles'
import { LogTableProps } from './types'

const LogTable = (props: LogTableProps) => {
  const [selectedRow, setSelectedRow] = useState<LogEntryModel | undefined>(
    undefined
  )

  return (
    <LogTableContainer className="ag-theme-balham">
      <AgGridReact<LogEntryModel>
        rowData={props.rowData}
        columnDefs={props.columnDefs}
        suppressFieldDotNotation
        enableCellTextSelection
        autoSizeStrategy={{ type: 'fitCellContents' }}
        onRowDoubleClicked={(e) => {
          setSelectedRow(e.data)
        }}
      />
      <LogTableDrawer
        anchor="right"
        open={selectedRow !== undefined}
        onClose={() => setSelectedRow(undefined)}
      >
        <LogTableDetailPre>
          {JSON.stringify(selectedRow, null, 2)}
        </LogTableDetailPre>
      </LogTableDrawer>
    </LogTableContainer>
  )
}

export default LogTable
