import Button from '@mui/material/Button'

import DeleteIcon from '@mui/icons-material/Delete'

import { CellRoot } from './styles'
import { TableActionsCellProps } from './types'

function TableActionsCell<T>({ data, onDelete }: TableActionsCellProps<T>) {
  return (
    <CellRoot>
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
    </CellRoot>
  )
}

export default TableActionsCell
