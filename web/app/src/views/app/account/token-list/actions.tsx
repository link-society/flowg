import { CustomCellRendererProps } from 'ag-grid-react'

import DeleteIcon from '@mui/icons-material/Delete'

import Button from '@mui/material/Button'

import { RowType } from './types'

type TokenActionsProps = CustomCellRendererProps<RowType> & {
  onDelete: (token: string) => void
}

export const TokenActions = ({ data, onDelete }: TokenActionsProps) => (
  <div className="h-full flex flex-row items-center justify-center">
    <Button
      variant="contained"
      size="small"
      color="error"
      onClick={() => onDelete(data!.token)}
    >
      <DeleteIcon />
    </Button>
  </div>
)
