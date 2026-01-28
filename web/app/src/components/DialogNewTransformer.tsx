import React, { useState } from 'react'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import TextField from '@mui/material/TextField'
import { DialogProps } from '@toolpad/core/useDialogs'

import CancelIcon from '@mui/icons-material/Cancel'
import SaveIcon from '@mui/icons-material/Save'

import * as configApi from '@/lib/api/operations/config'

import { useApiOperation } from '@/lib/hooks/api'

const DialogNewTransformer = ({
  open,
  onClose,
}: DialogProps<void, string | null>) => {
  const [name, setName] = useState('')

  const [onSubmit, loading] = useApiOperation(async () => {
    await configApi.saveTransformer(name, '')
    onClose(name)
  }, [name])

  return (
    <Dialog
      maxWidth="sm"
      fullWidth
      open={open}
      onClose={() => onClose(null)}
      slotProps={{
        paper: {
          component: 'form',
          onSubmit: (e: React.SubmitEvent<HTMLFormElement>) => {
            e.preventDefault()
            onSubmit()
          },
        },
      }}
    >
      <DialogTitle>Create new transformer</DialogTitle>
      <DialogContent>
        <div className="pt-3 w-full">
          <TextField
            id="input:transformers.modal.name"
            label="Transformer name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            type="text"
            variant="outlined"
            required
            className="w-full"
          />
        </div>
      </DialogContent>
      <DialogActions>
        <Button
          id="btn:transformers.modal.cancel"
          variant="contained"
          startIcon={<CancelIcon />}
          onClick={() => onClose(null)}
          disabled={loading}
        >
          Cancel
        </Button>
        <Button
          id="btn:transformers.modal.save"
          variant="contained"
          color="secondary"
          startIcon={!loading && <SaveIcon />}
          disabled={loading}
          type="submit"
        >
          {loading ? <CircularProgress color="inherit" size={24} /> : <>Save</>}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

export default DialogNewTransformer
