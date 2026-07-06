import { forwardRef, useImperativeHandle, useRef, useState } from 'react'

import { AgGridReact } from 'ag-grid-react'

import LogEntryModel from '@/lib/models/LogEntryModel'

import { LogTableContainer, LogTableDetailPre, LogTableDrawer } from './styles'
import { LogTableHandle, LogTableProps } from './types'

const LogTable = forwardRef<LogTableHandle, LogTableProps>((props, ref) => {
  const containerRef = useRef<HTMLDivElement>(null)
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

      const viewport = containerRef.current?.querySelector<HTMLElement>(
        '.ag-body-vertical-scroll-viewport'
      )
      const scrollTop = viewport?.scrollTop ?? 0
      const atTop = scrollTop <= 1
      const rowHeight = api.getDisplayedRowAtIndex(0)?.rowHeight ?? 0

      api.applyTransaction({ add: rows })

      if (viewport && !atTop && rowHeight > 0) {
        viewport.scrollTop = scrollTop + rows.length * rowHeight
      }
    },
  }))

  return (
    <LogTableContainer ref={containerRef} className="ag-theme-balham">
      <AgGridReact<LogEntryModel>
        ref={gridRef}
        rowData={props.rowData}
        columnDefs={props.columnDefs}
        suppressFieldDotNotation
        enableCellTextSelection
        animateRows={false}
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
