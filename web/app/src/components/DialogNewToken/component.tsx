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
>) => (
  <Dialog maxWidth="sm" fullWidth open={open} onClose={() => onClose()}>
    <DialogTitle>Your Personal Access Token</DialogTitle>
    <DialogContent>
      <FieldRow>
        <LabelIcon sx={{ mr: 1, my: 0.5 }} />
        <TextField
          id="input:account.tokens.modal.token_uuid"
          label="Token UUID"
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
        <KeyIcon sx={{ mr: 1, my: 0.5 }} />
        <TextField
          id="input:account.tokens.modal.token"
          label="Token"
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
        This token will be dislayed only once. Make sure to copy it before
        closing this dialog.
      </Typography>
    </DialogContent>
    <DialogActions>
      <Button
        id="btn:account.tokens.modal.done"
        variant="contained"
        color="secondary"
        onClick={() => onClose()}
      >
        Done
      </Button>
    </DialogActions>
  </Dialog>
)

export default DialogNewToken
