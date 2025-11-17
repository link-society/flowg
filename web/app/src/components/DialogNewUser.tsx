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

import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import CancelIcon from '@mui/icons-material/Cancel'
import LockIcon from '@mui/icons-material/Lock'
import SaveIcon from '@mui/icons-material/Save'

import * as aclApi from '@/lib/api/operations/acls'

import { useApiOperation } from '@/lib/hooks/api'

import UserModel from '@/lib/models/UserModel'

import InputTransferList from '@/components/InputTransferList'

const DialogNewUser = ({
  open,
  payload,
  onClose,
}: DialogProps<string[], UserModel | null>) => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [roles, setRoles] = useState<string[]>([])

  const [onSubmit, loading] = useApiOperation(async () => {
    const user = {
      name: username,
      roles,
    }

    await aclApi.saveUser(user, password)
    onClose(user)
  }, [username, password, roles, onClose])

  return (
    <Dialog
      maxWidth="sm"
      fullWidth
      open={open}
      onClose={() => onClose(null)}
      slotProps={{
        paper: {
          component: 'form',
          onSubmit: (e: React.FormEvent<HTMLFormElement>) => {
            e.preventDefault()
            onSubmit()
          },
        },
      }}
    >
      <DialogTitle>Create a new user</DialogTitle>
      <DialogContent>
        <Box className="flex flex-col items-stretch gap-3">
          <Box className="flex flex-row items-end">
            <AccountCircleIcon sx={{ mr: 1, my: 0.5 }} />
            <TextField
              id="input:admin.users.modal.username"
              label="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              type="text"
              variant="standard"
              className="grow"
              required
            />
          </Box>

          <Box className="flex flex-row items-end">
            <LockIcon sx={{ mr: 1, my: 0.5 }} />
            <TextField
              id="input:admin.users.modal.password"
              label="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              type="password"
              variant="standard"
              className="grow"
              required
            />
          </Box>

          <Box
            id="field:admin.users.modal.roles"
            className="flex flex-col items-stretch gap-2"
          >
            <span className="font-semibold">Roles:</span>
            <InputTransferList<string>
              choices={payload}
              getItemId={(v) => v}
              renderItem={(v) => (
                <Tooltip title={v} placement="right-start">
                  <Chip label={v} size="small" />
                </Tooltip>
              )}
              onChoiceUpdate={(choices) => setRoles([...choices])}
            />
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button
          id="btn:admin.users.modal.cancel"
          variant="contained"
          startIcon={<CancelIcon />}
          onClick={() => onClose(null)}
          disabled={loading}
        >
          Cancel
        </Button>
        <Button
          id="btn:admin.users.modal.submit"
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

export default DialogNewUser
