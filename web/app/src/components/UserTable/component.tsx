import { useCallback, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'

import Chip from '@mui/material/Chip'
import Typography from '@mui/material/Typography'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'

import { ColDef, GridSizeChangedEvent } from 'ag-grid-community'
import { AgGridReact, CustomCellRendererProps } from 'ag-grid-react'

import * as aclApi from '@/lib/api/operations/acls'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import UserModel from '@/lib/models/UserModel'

import ButtonNewUser from '@/components/ButtonNewUser/component'
import TableActionsCell from '@/components/TableActionsCell/component'

import {
  RolesCellRoot,
  UserTableCard,
  UserTableCardContent,
  UserTableCardHeader,
  UserTableCardHeaderTitle,
  UserTableCardHeaderTitleText,
} from './styles'
import { UserTableProps } from './types'

type RolesCellProps = CustomCellRendererProps<string[]>

const RolesCell = (props: RolesCellProps) => (
  <RolesCellRoot>
    {((props.value as string[] | null) ?? []).map((role) => (
      <Chip key={role} label={role} size="small" />
    ))}
  </RolesCellRoot>
)

const UserTable = ({ roles, users, defaultRoles }: UserTableProps) => {
  const { t } = useTranslation()
  const { permissions } = useProfile()
  const notify = useNotify()

  const gridRef = useRef<AgGridReact<UserModel>>(null!)

  const onGridSizeChanged = useCallback(
    (e: GridSizeChangedEvent<UserModel>) => {
      e.api.resetRowHeights()
    },
    []
  )

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

      notify.success(t('components.userTable.notifications.deleted'))
    },
    [gridRef]
  )

  const [rowData] = useState<UserModel[]>(users)
  const [columnDefs] = useState<ColDef<UserModel>[]>([
    {
      headerName: t('components.userTable.columns.username'),
      field: 'name',
      suppressMovable: true,
      sortable: false,
    },
    {
      headerName: t('components.userTable.columns.roles'),
      field: 'roles',
      cellDataType: false,
      cellRenderer: RolesCell,
      wrapText: true,
      autoHeight: true,
      suppressMovable: true,
      sortable: false,
      flex: 1,
    },
    {
      hide: !permissions.can_edit_acls,
      headerName: t('common.tableColumns.actions'),
      headerClass: 'flowg-actions-header',
      cellClass: 'flowg-actions-cell',
      cellRenderer: TableActionsCell,
      cellRendererParams: {
        onDelete: permissions.can_edit_acls ? onDelete : undefined,
      },
      suppressMovable: true,
      sortable: false,
    },
  ])

  return (
    <UserTableCard>
      <UserTableCardHeader
        title={
          <UserTableCardHeaderTitle>
            <AccountCircleIcon />
            <UserTableCardHeaderTitleText>
              <Typography variant="titleSm">
                {t('components.userTable.title')}
              </Typography>
            </UserTableCardHeaderTitleText>
            {permissions.can_edit_acls && (
              <ButtonNewUser
                roles={roles}
                defaultRoles={defaultRoles}
                onUserCreated={onNewUser}
              />
            )}
          </UserTableCardHeaderTitle>
        }
      />
      <UserTableCardContent id="table:admin.users">
        <AgGridReact<UserModel>
          ref={gridRef}
          loading={loading}
          rowData={rowData}
          columnDefs={columnDefs}
          enableCellTextSelection
          suppressHorizontalScroll
          getRowId={({ data }) => data.name}
          onGridSizeChanged={onGridSizeChanged}
          className="ag-theme-material flowg-table"
        />
      </UserTableCardContent>
    </UserTableCard>
  )
}

export default UserTable
