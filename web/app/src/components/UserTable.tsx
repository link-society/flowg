import { useCallback, useRef, useState } from 'react'

import Card from '@mui/material/Card'
import CardContent from '@mui/material/CardContent'
import CardHeader from '@mui/material/CardHeader'
import Chip from '@mui/material/Chip'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'

import { ColDef } from 'ag-grid-community'
import { AgGridReact, CustomCellRendererProps } from 'ag-grid-react'

import * as aclApi from '@/lib/api/operations/acls'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import UserModel from '@/lib/models/UserModel'

import ButtonNewUser from '@/components/ButtonNewUser'
import TableActionsCell from '@/components/TableActionsCell'

type RolesCellProps = CustomCellRendererProps<string[]>

const RolesCell = (props: RolesCellProps) => (
  <>
    {((props.value as string[] | null) ?? []).map((role) => (
      <Chip key={role} label={role} size="small" className="mx-1" />
    ))}
  </>
)

type UserTableProps = Readonly<{
  roles: string[]
  users: UserModel[]
}>

const UserTable = ({ roles, users }: UserTableProps) => {
  const { permissions } = useProfile()
  const notify = useNotify()

  const gridRef = useRef<AgGridReact<UserModel>>(null!)

  const onNewUser = useCallback(
    (user: UserModel) => {
      gridRef.current.api.applyTransaction({
        add: [user],
      })
    },
    [gridRef]
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
    [gridRef]
  )

  const [rowData] = useState<UserModel[]>(users)
  const [columnDefs] = useState<ColDef<UserModel>[]>([
    {
      hide: !permissions.can_edit_acls,
      headerName: 'Actions',
      headerClass: 'flowg-actions-header',
      cellRenderer: TableActionsCell,
      cellRendererParams: {
        onDelete: permissions.can_edit_acls ? onDelete : undefined,
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
            {permissions.can_edit_acls && (
              <ButtonNewUser roles={roles} onUserCreated={onNewUser} />
            )}
          </div>
        }
        className="bg-blue-400 text-white shadow-lg z-20"
      />
      <CardContent
        id="table:admin.users"
        className="p-0! grow shrink h-0 ag-theme-material flowg-table"
      >
        <AgGridReact<UserModel>
          ref={gridRef}
          loading={loading}
          rowData={rowData}
          columnDefs={columnDefs}
          enableCellTextSelection
          autoSizeStrategy={{ type: 'fitCellContents' }}
          getRowId={({ data }) => data.name}
        />
      </CardContent>
    </Card>
  )
}

export default UserTable
