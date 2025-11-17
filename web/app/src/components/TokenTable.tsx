import { useCallback, useRef, useState } from 'react'

import Card from '@mui/material/Card'
import CardContent from '@mui/material/CardContent'
import CardHeader from '@mui/material/CardHeader'

import KeyIcon from '@mui/icons-material/Key'

import { ColDef } from 'ag-grid-community'
import { AgGridReact, CustomCellRendererProps } from 'ag-grid-react'

import * as tokenApi from '@/lib/api/operations/token'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import ButtonNewToken from '@/components/ButtonNewToken'
import TableActionsCell from '@/components/TableActionsCell'

type RowType = { token: string }

type TokenCellProps = CustomCellRendererProps<string>

const TokenCell = (props: TokenCellProps) => (
  <span className="font-mono">{props.value}</span>
)

type TokenTableProps = Readonly<{
  tokens: string[]
}>

const TokenTable = ({ tokens }: TokenTableProps) => {
  const notify = useNotify()

  const gridRef = useRef<AgGridReact<RowType>>(null!)

  const onNewToken = useCallback(
    (token: string) => {
      gridRef.current.api.applyTransaction({
        add: [{ token }],
      })
    },
    [gridRef]
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
    [gridRef]
  )

  const [rowData] = useState<RowType[]>(tokens.map((token) => ({ token })))
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
      cellRenderer: TableActionsCell,
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
            <ButtonNewToken onTokenCreated={onNewToken} />
          </div>
        }
        className="bg-blue-400 text-white shadow-lg z-20"
      />
      <CardContent
        id="table:account.tokens"
        className="p-0! grow shrink h-0 ag-theme-material flowg-table"
      >
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

export default TokenTable
