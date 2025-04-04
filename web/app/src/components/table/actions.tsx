import Button from '@mui/material/Button'

import DeleteIcon from '@mui/icons-material/Delete'

import { CustomCellRendererProps } from 'ag-grid-react'

type ActionsProps<T> = CustomCellRendererProps<T> & {
  onDelete?: (data: T) => void
}

export function Actions<T>({ data, onDelete }: ActionsProps<T>) {
  return (
    <div className="h-full flex flex-row items-center justify-center">
      {onDelete && (
        <Button
          variant="contained"
          size="small"
          color="error"
          onClick={() => onDelete(data!)}
          data-ref="btn:generic.tablerow.actions.delete"
        >
          <DeleteIcon />
        </Button>
      )}
    </div>
  )
}
