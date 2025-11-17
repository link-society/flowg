import { useCallback, useRef, useState } from 'react'

import Card from '@mui/material/Card'
import CardContent from '@mui/material/CardContent'
import CardHeader from '@mui/material/CardHeader'
import Chip from '@mui/material/Chip'

import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings'

import { ColDef } from 'ag-grid-community'
import { AgGridReact, CustomCellRendererProps } from 'ag-grid-react'

import * as aclApi from '@/lib/api/operations/acls'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import RoleModel from '@/lib/models/RoleModel'
import { ScopeLabels } from '@/lib/models/Scopes'

import ButtonNewRole from '@/components/ButtonNewRole'
import TableActionsCell from '@/components/TableActionsCell'

type ScopesCellProps = CustomCellRendererProps<string[]>

const ScopesCell = (props: ScopesCellProps) => (
  <>
    {((props.value as string[] | null) ?? []).map((scope) => (
      <Chip
        key={scope}
        label={ScopeLabels[scope as keyof typeof ScopeLabels] ?? '#ERR#'}
        size="small"
        className="mx-1"
      />
    ))}
  </>
)

type RoleTableProps = Readonly<{
  roles: RoleModel[]
}>

const RoleTable = ({ roles }: RoleTableProps) => {
  const { permissions } = useProfile()
  const notify = useNotify()

  const gridRef = useRef<AgGridReact<RoleModel>>(null!)

  const onNewRole = useCallback(
    (role: RoleModel) => {
      gridRef.current.api.applyTransaction({
        add: [role],
      })
    },
    [gridRef]
  )

  const [onDelete, loading] = useApiOperation(
    async (role: RoleModel) => {
      await aclApi.deleteRole(role.name)

      const rowNode = gridRef.current.api.getRowNode(role.name)
      if (rowNode !== undefined && rowNode.data !== undefined) {
        gridRef.current.api.applyTransaction({
          remove: [rowNode.data],
        })
      }

      notify.success('Role deleted')
    },
    [gridRef]
  )

  const [rowData] = useState<RoleModel[]>(roles)
  const [columnDefs] = useState<ColDef<RoleModel>[]>([
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

  return (
    <Card className="max-lg:min-h-96 lg:h-full flex flex-col items-stretch">
      <CardHeader
        title={
          <div className="flex items-center gap-3">
            <AdminPanelSettingsIcon />
            <span className="grow">Roles</span>
            {permissions.can_edit_acls && (
              <ButtonNewRole onRoleCreated={onNewRole} />
            )}
          </div>
        }
        className="bg-blue-400 text-white shadow-lg z-20"
      />
      <CardContent
        id="table:admin.roles"
        className="p-0! grow shrink h-0 ag-theme-material flowg-table"
      >
        <AgGridReact<RoleModel>
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

export default RoleTable
