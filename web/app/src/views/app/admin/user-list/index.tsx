import { useCallback, useRef, useState } from 'react'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'

import Card from '@mui/material/Card'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'

import { AgGridReact } from 'ag-grid-react'
import { ColDef } from 'ag-grid-community'

import { Actions } from '@/components/table/actions'
import { CreateUserButton } from './create-btn'
import { RolesCell } from './roles-cell'

import * as aclApi from '@/lib/api/operations/acls'
import { UserModel } from '@/lib/models'

type UserListProps = Readonly<{
  roles: string[]
  users: UserModel[]
}>

export const UserList = ({ roles, users }: UserListProps) => {
  const notify = useNotify()

  const gridRef = useRef<AgGridReact<UserModel>>(null!)

  const onNewUser = useCallback(
    (user: UserModel) => {
      gridRef.current.api.applyTransaction({
        add: [user],
      })
    },
    [gridRef],
  )

  const [onDelete, loading] = useApiOperation(
    async (user: UserModel) => {
      await aclApi.deleteUser(user.name)

      const rowNode = gridRef.current.api.getRowNode(user.name)
      if (rowNode !== undefined && rowNode.data !== undefined) {
        gridRef.current.api.applyTransaction({
          remove: [rowNode.data],
        })
      }

      notify.success('User deleted')
    },
    [gridRef],
  )

  const [rowData] = useState<UserModel[]>(users)
  const [columnDefs] = useState<ColDef<UserModel>[]>([
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

  return (
    <Card className="max-lg:min-h-96 lg:h-full flex flex-col items-stretch">
      <CardHeader
        title={
          <div className="flex items-center gap-3">
            <AccountCircleIcon />
            <span className="grow">Users</span>
            <CreateUserButton roles={roles} onUserCreated={onNewUser} />
          </div>
        }
        className="bg-blue-400 text-white shadow-lg z-20"
      />
      <CardContent className="p-0! grow shrink h-0 ag-theme-material flowg-table">
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
  )
}
