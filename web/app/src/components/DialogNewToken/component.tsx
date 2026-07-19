import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'

import KeyIcon from '@mui/icons-material/Key'
import LabelIcon from '@mui/icons-material/Label'

import { DialogProps } from '@/lib/models/Dialog'

import { FieldRow } from './styles'

const DialogNewToken = ({
  open,
  payload,
  onClose,
}: DialogProps<
  {
    token: string
    token_uuid: string
  },
  void
>) => {
  const { t } = useTranslation()

  return (
    <Dialog maxWidth="sm" fullWidth open={open} onClose={() => onClose()}>
      <DialogTitle>{t('components.dialogNewToken.title')}</DialogTitle>
      <DialogContent>
        <FieldRow>
          <LabelIcon />
          <TextField
            id="input:account.tokens.modal.token_uuid"
            label={t('components.dialogNewToken.tokenUuidLabel')}
            type="text"
            value={payload.token_uuid}
            variant="standard"
            fullWidth
            slotProps={{
              input: {
                readOnly: true,
              },
            }}
          />
        </FieldRow>

        <FieldRow>
          <KeyIcon />
          <TextField
            id="input:account.tokens.modal.token"
            label={t('components.dialogNewToken.tokenLabel')}
            value={payload.token}
            type="text"
            variant="standard"
            fullWidth
            slotProps={{
              input: {
                readOnly: true,
              },
            }}
          />
        </FieldRow>

        <Typography variant="text">
          {t('components.dialogNewToken.hint')}
        </Typography>
      </DialogContent>
      <DialogActions>
        <Button
          id="btn:account.tokens.modal.done"
          variant="contained"
          color="secondary"
          onClick={() => onClose()}
        >
          {t('components.dialogNewToken.done')}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

export default DialogNewToken
