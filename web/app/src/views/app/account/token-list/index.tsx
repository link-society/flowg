import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useSnackbar } from 'notistack'

import KeyIcon from '@mui/icons-material/Key'
import AddIcon from '@mui/icons-material/Add'

import Card from '@mui/material/Card'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'
import Button from '@mui/material/Button'
import Modal from '@mui/material/Modal'

import { AgGridReact } from 'ag-grid-react'
import { ColDef } from 'ag-grid-community'

import { RowType } from './types'
import { TokenCell } from './cell'
import { TokenActions } from './actions'

import './style.css'

import { UnauthenticatedError } from '@/lib/api/errors'
import * as tokenApi from '@/lib/api/operations/token'

type TokenListProps = {
  tokens: string[]
}

export const TokenList = ({ tokens }: TokenListProps) => {
  const navigate = useNavigate()
  const { enqueueSnackbar } = useSnackbar()

  const [loading, setLoading] = useState(false)
  const [rowData, setRowData] = useState<RowType[]>(
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
      headerClass: 'flowg-pat-actions-header',
      cellRenderer: TokenActions,
      cellRendererParams: {
        onDelete: async (token: string) => {
          setLoading(true)

          try {
            await tokenApi.deleteToken(token)
            setRowData(rowData.filter(row => row.token !== token))

            enqueueSnackbar('Token deleted', { variant: 'success' })
          }
          catch (error) {
            if (error instanceof UnauthenticatedError) {
              enqueueSnackbar('Session expired', { variant: 'error' })
              navigate('/web/login')
            }
            else {
              enqueueSnackbar('Unknown error', { variant: 'error' })
            }

            console.error(error)
          }

          setLoading(false)
        },
      },
      suppressMovable: true,
      sortable: false,
    },
  ])

  return (
    <>
      <Card className="lg:h-full lg:flex lg:flex-col lg:items-stretch">
        <CardHeader
          title={
            <div className="flex items-center gap-3">
              <KeyIcon />
              <span className="flex-grow">API Tokens</span>
              <Button
                variant="contained"
                color="secondary"
                size="small"
                startIcon={<AddIcon />}
              >
                New Token
              </Button>
            </div>
          }
          className="bg-blue-400 text-white shadow-lg z-20"
        />
        <CardContent className="!p-0 lg:flex-grow ag-theme-material flowg-pat-table">
          <AgGridReact<RowType>
            loading={loading}
            rowData={rowData}
            columnDefs={columnDefs}
            enableCellTextSelection
          />
        </CardContent>
      </Card>

      <Modal></Modal>
    </>
  )
}
