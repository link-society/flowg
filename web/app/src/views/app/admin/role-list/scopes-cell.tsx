import Chip from '@mui/material/Chip'

import { CustomCellRendererProps } from 'ag-grid-react'

import { SCOPE_LABELS } from '@/lib/models/permissions'

type ScopesCellProps = CustomCellRendererProps<string[]> & {}

export const ScopesCell = (props: ScopesCellProps) => (
  <>
    {(props.value as string[]).map((scope, idx) => (
      <Chip
        key={idx}
        label={SCOPE_LABELS[scope as keyof typeof SCOPE_LABELS] ?? '#ERR#'}
        size="small"
        className="mx-1"
      />
    ))}
  </>
)
