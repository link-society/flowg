import Chip from '@mui/material/Chip'

import { CustomCellRendererProps } from 'ag-grid-react'

type RolesCellProps = CustomCellRendererProps<string[]>

export const RolesCell = (props: RolesCellProps) => (
  <>
    {(props.value as string[]).map((role) => (
      <Chip
        key={role}
        label={role}
        size="small"
        className="mx-1"
      />
    ))}
  </>
)
