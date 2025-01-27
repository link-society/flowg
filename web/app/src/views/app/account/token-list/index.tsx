import { useCallback, useRef, useState } from 'react'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import KeyIcon from '@mui/icons-material/Key'

import Card from '@mui/material/Card'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'

import { AgGridReact } from 'ag-grid-react'
import { ColDef } from 'ag-grid-community'

import { Actions } from '@/components/table/actions'

import { CreateTokenButton } from './create-btn'
import { RowType } from './types'
import { TokenCell } from './cell'

import * as tokenApi from '@/lib/api/operations/token'

type TokenListProps = Readonly<{
  tokens: string[]
}>

export const TokenList = ({ tokens }: TokenListProps) => {
  const notify = useNotify()

  const gridRef = useRef<AgGridReact<RowType>>(null!)

  const onNewToken = useCallback(
    (token: string) => {
      gridRef.current.api.applyTransaction({
        add: [{ token }],
      })
    },
    [gridRef],
  )

  const [onDelete, loading] = useApiOperation(
    async (data: RowType) => {
      await tokenApi.deleteToken(data.token)

      const rowNode = gridRef.current.api.getRowNode(data.token)
      if (rowNode !== undefined && rowNode.data !== undefined) {
        gridRef.current.api.applyTransaction({
          remove: [rowNode.data],
        })
      }

      notify.success('Token deleted')
    },
    [gridRef],
  )

  const [rowData] = useState<RowType[]>(
    tokens.map(token => ({ token })),
  )
  const [columnDefs] = useState<ColDef<RowType>[]>([
    {
      headerName: 'Token',
      field: 'token',
      cellRenderer: TokenCell,
      suppressMovable: true,
      sortable: false,
      flex: 1,
    },
    {
      headerName: 'Actions',
      headerClass: 'flowg-actions-header',
      cellRenderer: Actions,
      cellRendererParams: {
        onDelete,
      },
      suppressMovable: true,
      sortable: false,
    },
  ])

  return (
    <Card className="max-lg:min-h-96 lg:h-full flex flex-col items-stretch">
      <CardHeader
        title={
          <div className="flex items-center gap-3">
            <KeyIcon />
            <span className="grow">API Tokens</span>
            <CreateTokenButton onTokenCreated={onNewToken} />
          </div>
        }
        className="bg-blue-400 text-white shadow-lg z-20"
      />
      <CardContent className="p-0! grow shrink h-0 ag-theme-material flowg-table">
        <AgGridReact<RowType>
          ref={gridRef}
          loading={loading}
          rowData={rowData}
          columnDefs={columnDefs}
          enableCellTextSelection
          getRowId={({ data }) => data.token}
        />
      </CardContent>
    </Card>
  )
}
