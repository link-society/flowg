import { useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'

import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings'

import Card from '@mui/material/Card'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'

import { AgGridReact } from 'ag-grid-react'
import { ColDef } from 'ag-grid-community'

import { Actions } from '@/components/table/actions'
import { CreateRoleButton } from './create-btn'
import { ScopesCell } from './scopes-cell'

import { UnauthenticatedError } from '@/lib/api/errors'
import * as aclApi from '@/lib/api/operations/acls'
import { RoleModel } from '@/lib/models'

type RoleListProps = {
  roles: RoleModel[]
}

export const RoleList = ({ roles }: RoleListProps) => {
  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()

  const [loading, setLoading] = useState(false)

  const gridRef = useRef<AgGridReact<RoleModel>>(null!)
  const [rowData] = useState<RoleModel[]>(roles)
  const [columnDefs] = useState<ColDef<RoleModel>[]>([

    {
      headerName: 'Actions',
      headerClass: 'flowg-actions-header',
      cellRenderer: Actions,
      cellRendererParams: {
        onDelete: async (role: RoleModel) => {
          setLoading(true)

          try {
            await aclApi.deleteRole(role.name)

            const rowNode = gridRef.current.api.getRowNode(role.name)
            if (rowNode !== undefined) {
              gridRef.current.api.applyTransaction({
                remove: [rowNode.data],
              })
            }

            notifications.show('Role deleted', {
              severity: 'success',
              autoHideDuration: config.notifications?.autoHideDuration,
            })
          }
          catch (error) {
            if (error instanceof UnauthenticatedError) {
              notifications.show('Session expired', {
                severity: 'error',
                autoHideDuration: config.notifications?.autoHideDuration,
              })
              navigate('/web/login')
            }
            else {
              notifications.show('Unknown error', {
                severity: 'error',
                autoHideDuration: config.notifications?.autoHideDuration,
              })
            }

            console.error(error)
          }

          setLoading(false)
        },
      },
      suppressMovable: true,
      sortable: false,
    },
    {
      headerName: 'Role',
      field: 'name',
      suppressMovable: true,
      sortable: false,
    },
    {
      headerName: 'Permissions',
      field: 'scopes',
      cellDataType: false,
      cellRenderer: ScopesCell,
      suppressMovable: true,
      sortable: false,
    },
  ])

  const onNewRole = (role: RoleModel) => {
    gridRef.current.api.applyTransaction({
      add: [role],
    })
  }

  return (
    <>
      <Card className="max-lg:min-h-96 lg:h-full flex flex-col items-stretch">
        <CardHeader
          title={
            <div className="flex items-center gap-3">
              <AdminPanelSettingsIcon />
              <span className="flex-grow">Roles</span>
              <CreateRoleButton onRoleCreated={onNewRole} />
            </div>
          }
          className="bg-blue-400 text-white shadow-lg z-20"
        />
        <CardContent className="!p-0 flex-grow flex-shrink h-0 ag-theme-material flowg-table">
          <AgGridReact<RoleModel>
            ref={gridRef}
            loading={loading}
            rowData={rowData}
            columnDefs={columnDefs}
            enableCellTextSelection
            autoSizeStrategy={{type: 'fitCellContents'}}
            getRowId={({ data }) => data.name}
          />
        </CardContent>
      </Card>
    </>
  )
}
