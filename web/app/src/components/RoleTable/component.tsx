import { useCallback, useRef, useState } from 'react'

import Chip from '@mui/material/Chip'
import Typography from '@mui/material/Typography'

import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings'

import { ColDef, GridSizeChangedEvent } from 'ag-grid-community'
import { AgGridReact, CustomCellRendererProps } from 'ag-grid-react'

import * as aclApi from '@/lib/api/operations/acls'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'
import { useProfile } from '@/lib/hooks/profile'

import RoleModel from '@/lib/models/RoleModel'
import { ScopeLabels } from '@/lib/models/Scopes'

import ButtonNewRole from '@/components/ButtonNewRole/component'
import TableActionsCell from '@/components/TableActionsCell/component'

import {
  RoleTableCard,
  RoleTableCardContent,
  RoleTableCardHeader,
  RoleTableCardHeaderTitle,
  RoleTableCardHeaderTitleText,
  ScopesCellRoot,
} from './styles'
import { RoleTableProps } from './types'

type ScopesCellProps = CustomCellRendererProps<string[]>

const ScopesCell = (props: ScopesCellProps) => (
  <ScopesCellRoot>
    {((props.value as string[] | null) ?? []).map((scope) => (
      <Chip
        key={scope}
        label={ScopeLabels[scope as keyof typeof ScopeLabels] ?? '#ERR#'}
        size="small"
      />
    ))}
  </ScopesCellRoot>
)

const RoleTable = ({ roles }: RoleTableProps) => {
  const { permissions } = useProfile()
  const notify = useNotify()

  const gridRef = useRef<AgGridReact<RoleModel>>(null!)

  const onGridSizeChanged = useCallback(
    (e: GridSizeChangedEvent<RoleModel>) => {
      e.api.resetRowHeights()
    },
    []
  )

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
      wrapText: true,
      autoHeight: true,
      suppressMovable: true,
      sortable: false,
      flex: 1,
    },
    {
      hide: !permissions.can_edit_acls,
      headerName: 'Actions',
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
    <RoleTableCard>
      <RoleTableCardHeader
        title={
          <RoleTableCardHeaderTitle>
            <AdminPanelSettingsIcon />
            <RoleTableCardHeaderTitleText>
              <Typography variant="titleSm">Roles</Typography>
            </RoleTableCardHeaderTitleText>
            {permissions.can_edit_acls && (
              <ButtonNewRole onRoleCreated={onNewRole} />
            )}
          </RoleTableCardHeaderTitle>
        }
      />
      <RoleTableCardContent id="table:admin.roles">
        <AgGridReact<RoleModel>
          ref={gridRef}
          loading={loading}
          rowData={rowData}
          columnDefs={columnDefs}
          enableCellTextSelection
          getRowId={({ data }) => data.name}
          onGridSizeChanged={onGridSizeChanged}
          className="ag-theme-material flowg-table"
        />
      </RoleTableCardContent>
    </RoleTableCard>
  )
}

export default RoleTable
