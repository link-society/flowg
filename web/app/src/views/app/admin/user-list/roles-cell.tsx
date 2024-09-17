import Chip from '@mui/material/Chip'

import { CustomCellRendererProps } from 'ag-grid-react'

type RolesCellProps = CustomCellRendererProps<string[]> & {}

export const RolesCell = (props: RolesCellProps) => (
  <>
    {(props.value as string[]).map((role, idx) => (
      <Chip
        key={idx}
        label={role}
        size="small"
        className="mx-1"
      />
    ))}
  </>
)
