import { forwardRef, useImperativeHandle, useRef, useState } from 'react'

import { AgGridReact } from 'ag-grid-react'

import LogEntryModel from '@/lib/models/LogEntryModel'

import { LogTableContainer, LogTableDetailPre, LogTableDrawer } from './styles'
import { LogTableHandle, LogTableProps } from './types'

const LogTable = forwardRef<LogTableHandle, LogTableProps>((props, ref) => {
  const gridRef = useRef<AgGridReact<LogEntryModel>>(null)
  const [selectedRow, setSelectedRow] = useState<LogEntryModel | undefined>(
    undefined
  )

  useImperativeHandle(ref, () => ({
    appendRows: (rows) => {
      gridRef.current?.api?.applyTransaction({ add: rows })
    },
  }))

  return (
    <LogTableContainer className="ag-theme-balham">
      <AgGridReact<LogEntryModel>
        ref={gridRef}
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
})

export default LogTable
