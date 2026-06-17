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
      const api = gridRef.current?.api
      if (!api) {
        return
      }

      const scrollTop = api.getVerticalPixelRange().top
      const atTop = scrollTop <= 1

      api.applyTransaction({ add: rows })

      if (!atTop) {
        const rowHeight = api.getDisplayedRowAtIndex(0)?.rowHeight ?? 0
        if (rowHeight > 0) {
          const firstVisibleIndex = Math.round(scrollTop / rowHeight)
          api.ensureIndexVisible(firstVisibleIndex + rows.length, 'top')
        }
      }
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
