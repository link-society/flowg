import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import TextField from '@mui/material/TextField'
import { DialogProps } from '@toolpad/core/useDialogs'

import KeyIcon from '@mui/icons-material/Key'
import LabelIcon from '@mui/icons-material/Label'

export const ShowNewTokenModal = ({
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
      <Box className="flex flex-row items-end mb-3">
        <LabelIcon sx={{ mr: 1, my: 0.5 }} />
        <TextField
          id="input:account.tokens.modal.token_uuid"
          label="Token UUID"
          type="text"
          value={payload.token_uuid}
          variant="standard"
          className="grow"
          slotProps={{
            input: {
              readOnly: true,
            },
          }}
        />
      </Box>

      <Box className="flex flex-row items-end mb-3">
        <KeyIcon sx={{ mr: 1, my: 0.5 }} />
        <TextField
          id="input:account.tokens.modal.token"
          label="Token"
          value={payload.token}
          type="text"
          variant="standard"
          className="grow"
          slotProps={{
            input: {
              readOnly: true,
            },
          }}
        />
      </Box>

      <p>
        This token will be dislayed only once. Make sure to copy it before
        closing this dialog.
      </p>
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
