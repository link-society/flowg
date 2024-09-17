import { useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'

import Card from '@mui/material/Card'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'

import { AgGridReact } from 'ag-grid-react'
import { ColDef } from 'ag-grid-community'

import { Actions } from '@/components/table/actions'
import { CreateUserButton } from './create-btn'
import { RolesCell } from './roles-cell'

import { UnauthenticatedError } from '@/lib/api/errors'
import * as aclApi from '@/lib/api/operations/acls'
import { UserModel } from '@/lib/models'

type UserListProps = {
  roles: string[]
  users: UserModel[]
}

export const UserList = ({ roles, users }: UserListProps) => {
  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()

  const [loading, setLoading] = useState(false)

  const gridRef = useRef<AgGridReact<UserModel>>(null!)
  const [rowData] = useState<UserModel[]>(users)
  const [columnDefs] = useState<ColDef<UserModel>[]>([
    {
      headerName: 'Actions',
      headerClass: 'flowg-actions-header',
      cellRenderer: Actions,
      cellRendererParams: {
        onDelete: async (user: UserModel) => {
          setLoading(true)

          try {
            await aclApi.deleteUser(user.name)

            const rowNode = gridRef.current.api.getRowNode(user.name)
            if (rowNode !== undefined) {
              gridRef.current.api.applyTransaction({
                remove: [rowNode.data],
              })
            }

            notifications.show('User deleted', {
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
      headerName: 'Username',
      field: 'name',
      suppressMovable: true,
      sortable: false,
    },
    {
      headerName: 'Roles',
      field: 'roles',
      cellDataType: false,
      cellRenderer: RolesCell,
      suppressMovable: true,
      sortable: false,
    },
  ])

  const onNewUser = (user: UserModel) => {
    gridRef.current.api.applyTransaction({
      add: [user],
    })
  }

  return (
    <>
      <Card className="max-lg:min-h-96 lg:h-full flex flex-col items-stretch">
        <CardHeader
          title={
            <div className="flex items-center gap-3">
              <AccountCircleIcon />
              <span className="flex-grow">Users</span>
              <CreateUserButton roles={roles} onUserCreated={onNewUser} />
            </div>
          }
          className="bg-blue-400 text-white shadow-lg z-20"
        />
        <CardContent className="!p-0 flex-grow flex-shrink h-0 ag-theme-material flowg-table">
          <AgGridReact<UserModel>
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
