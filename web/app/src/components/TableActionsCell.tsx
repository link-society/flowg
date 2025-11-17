import Button from '@mui/material/Button'

import DeleteIcon from '@mui/icons-material/Delete'

import { CustomCellRendererProps } from 'ag-grid-react'

type TableActionsCellProps<T> = CustomCellRendererProps<T> & {
  onDelete?: (data: T) => void
}

function TableActionsCell<T>({ data, onDelete }: TableActionsCellProps<T>) {
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

export default TableActionsCell
