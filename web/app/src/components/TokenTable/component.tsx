import { useCallback, useRef, useState } from 'react'

import Typography from '@mui/material/Typography'

import KeyIcon from '@mui/icons-material/Key'

import { ColDef } from 'ag-grid-community'
import { AgGridReact, CustomCellRendererProps } from 'ag-grid-react'

import * as tokenApi from '@/lib/api/operations/token'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import ButtonNewToken from '@/components/ButtonNewToken/component'
import TableActionsCell from '@/components/TableActionsCell/component'

import {
  TokenCellRoot,
  TokenTableCard,
  TokenTableCardContent,
  TokenTableCardHeader,
  TokenTableCardHeaderTitle,
  TokenTableCardHeaderTitleText,
} from './styles'
import { RowType, TokenTableProps } from './types'

type TokenCellProps = CustomCellRendererProps<string>

const TokenCell = (props: TokenCellProps) => (
  <TokenCellRoot>{props.value}</TokenCellRoot>
)

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
    <TokenTableCard>
      <TokenTableCardHeader
        title={
          <TokenTableCardHeaderTitle>
            <KeyIcon />
            <TokenTableCardHeaderTitleText>
              <Typography variant="titleSm">API Tokens</Typography>
            </TokenTableCardHeaderTitleText>
            <ButtonNewToken onTokenCreated={onNewToken} />
          </TokenTableCardHeaderTitle>
        }
      />
      <TokenTableCardContent
        id="table:account.tokens"
        className="ag-theme-material flowg-table"
      >
        <AgGridReact<RowType>
          ref={gridRef}
          loading={loading}
          rowData={rowData}
          columnDefs={columnDefs}
          enableCellTextSelection
          getRowId={({ data }) => data.token}
        />
      </TokenTableCardContent>
    </TokenTableCard>
  )
}

export default TokenTable
