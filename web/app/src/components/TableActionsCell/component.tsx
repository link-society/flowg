import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'

import DeleteIcon from '@mui/icons-material/Delete'

import { useDialogs } from '@/lib/hooks/dialogs'

import DialogConfirm from '@/components/DialogConfirm/component'

import { CellRoot } from './styles'
import { TableActionsCellProps } from './types'

function TableActionsCell<T>({ data, onDelete }: TableActionsCellProps<T>) {
  const { t } = useTranslation()
  const dialogs = useDialogs()

  const handleDeleteClick = async () => {
    if (!onDelete || !data) return

    const confirmed = await dialogs.open(DialogConfirm, {
      title: t('components.tableActionsCell.deleteConfirm.title'),
      message: t('components.tableActionsCell.deleteConfirm.message'),
      confirmLabel: t('common.actions.delete'),
      danger: true,
    })

    if (confirmed) {
      onDelete(data)
    }
  }

  return (
    <CellRoot>
      {onDelete && data && (
        <Button
          variant="contained"
          size="small"
          color="error"
          onClick={handleDeleteClick}
          data-ref="btn:generic.tablerow.actions.delete"
        >
          <DeleteIcon />
        </Button>
      )}
    </CellRoot>
  )
}

export default TableActionsCell
