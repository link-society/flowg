import KeyIcon from '@mui/icons-material/Key'

import Dialog from '@mui/material/Dialog'
import DialogTitle from '@mui/material/DialogTitle'
import DialogContent from '@mui/material/DialogContent'
import DialogActions from '@mui/material/DialogActions'
import Box from '@mui/material/Box'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'

import { DialogProps } from '@toolpad/core/useDialogs'

export const ShowNewTokenModal = ({ open, payload, onClose }: DialogProps<string, void>) => (
  <Dialog
    maxWidth="sm"
    fullWidth
    open={open}
    onClose={() => onClose()}
  >
    <DialogTitle>Your Personal Access Token</DialogTitle>
    <DialogContent>
      <Box className="flex flex-row items-end">
        <KeyIcon sx={{ mr: 1, my: 0.5 }} />
        <TextField
          label="Token"
          value={payload}
          type="text"
          variant="standard"
          className="flex-grow"
          slotProps={{
            input: {
              readOnly: true,
            },
          }}
        />
      </Box>

      <p className="mt-3">
        This token will be dislayed only once.
        Make sure to copy it before closing this dialog.
      </p>
    </DialogContent>
    <DialogActions>
      <Button
        variant="contained"
        color="secondary"
        onClick={() => onClose()}
      >
        Done
      </Button>
    </DialogActions>
  </Dialog>
)
