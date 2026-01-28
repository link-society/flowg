import React, { useState } from 'react'

import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import Chip from '@mui/material/Chip'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import TextField from '@mui/material/TextField'
import Tooltip from '@mui/material/Tooltip'
import { DialogProps } from '@toolpad/core/useDialogs'

import CancelIcon from '@mui/icons-material/Cancel'
import SaveIcon from '@mui/icons-material/Save'

import * as aclApi from '@/lib/api/operations/acls'

import { useApiOperation } from '@/lib/hooks/api'

import RoleModel from '@/lib/models/RoleModel'
import { ScopeLabels, Scopes } from '@/lib/models/Scopes'

import InputTransferList from '@/components/InputTransferList'

const DialogNewRole = ({
  open,
  onClose,
}: DialogProps<void, RoleModel | null>) => {
  const [name, setName] = useState('')
  const [scopes, setScopes] = useState<string[]>([])

  const [onSubmit, loading] = useApiOperation(async () => {
    const role = {
      name,
      scopes,
    }

    await aclApi.saveRole(role)
    onClose(role)
  }, [name, scopes, onClose])

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
      <DialogTitle>Create a new role</DialogTitle>
      <DialogContent>
        <Box className="flex flex-col items-stretch gap-3">
          <TextField
            id="input:admin.roles.modal.name"
            label="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            type="text"
            variant="outlined"
            className="mt-2!"
            required
          />

          <Box
            id="field:admin.roles.modal.scopes"
            className="flex flex-col items-stretch gap-2"
          >
            <span className="font-semibold">Permissions:</span>
            <InputTransferList<string>
              choices={Scopes}
              getItemId={(v) => v}
              renderItem={(v) => (
                <Tooltip
                  title={ScopeLabels[v as keyof typeof ScopeLabels] ?? '#ERR#'}
                  placement="right-start"
                >
                  <Chip
                    label={
                      ScopeLabels[v as keyof typeof ScopeLabels] ?? '#ERR#'
                    }
                    size="small"
                  />
                </Tooltip>
              )}
              onChoiceUpdate={(choices) => setScopes([...choices])}
            />
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button
          id="btn:admin.roles.modal.cancel"
          variant="contained"
          startIcon={<CancelIcon />}
          onClick={() => onClose(null)}
          disabled={loading}
        >
          Cancel
        </Button>
        <Button
          id="btn:admin.roles.modal.submit"
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

export default DialogNewRole
