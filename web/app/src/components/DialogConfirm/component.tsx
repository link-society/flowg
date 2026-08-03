import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogContentText from '@mui/material/DialogContentText'
import DialogTitle from '@mui/material/DialogTitle'

import CancelIcon from '@mui/icons-material/Cancel'
import CheckIcon from '@mui/icons-material/Check'
import DeleteIcon from '@mui/icons-material/Delete'

import { DialogProps } from '@/lib/models/Dialog'

import { DialogConfirmPayload } from './types'

const DialogConfirm = ({
  open,
  payload,
  onClose,
}: DialogProps<DialogConfirmPayload, boolean>) => {
  const { t } = useTranslation()
  const { title, message, confirmLabel, danger = false } = payload

  return (
    <Dialog open={open} onClose={() => onClose(false)} fullWidth maxWidth="xs">
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <DialogContentText>{message}</DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button
          id="btn:confirm-dialog.cancel"
          variant="contained"
          startIcon={<CancelIcon />}
          onClick={() => onClose(false)}
        >
          {t('common.actions.cancel')}
        </Button>
        <Button
          id="btn:confirm-dialog.confirm"
          variant="contained"
          color={danger ? 'error' : 'secondary'}
          startIcon={danger ? <DeleteIcon /> : <CheckIcon />}
          onClick={() => onClose(true)}
          autoFocus
        >
          {confirmLabel ?? t('common.actions.confirm')}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

export default DialogConfirm
