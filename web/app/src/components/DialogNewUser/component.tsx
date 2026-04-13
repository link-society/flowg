import React, { useState } from 'react'

import Button from '@mui/material/Button'
import Chip from '@mui/material/Chip'
import CircularProgress from '@mui/material/CircularProgress'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import TextField from '@mui/material/TextField'
import Tooltip from '@mui/material/Tooltip'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import CancelIcon from '@mui/icons-material/Cancel'
import LockIcon from '@mui/icons-material/Lock'
import SaveIcon from '@mui/icons-material/Save'

import * as aclApi from '@/lib/api/operations/acls'

import { useApiOperation } from '@/lib/hooks/api'

import { DialogProps } from '@/lib/models/Dialog'
import UserModel from '@/lib/models/UserModel'

import InputTransferList from '@/components/InputTransferList/component'

import { FieldLabel, FieldRow, FieldStack, FormStack } from './styles'

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
          onSubmit: (e: React.SubmitEvent<HTMLFormElement>) => {
            e.preventDefault()
            onSubmit()
          },
        },
      }}
    >
      <DialogTitle>Create a new user</DialogTitle>
      <DialogContent>
        <FormStack>
          <FieldRow>
            <AccountCircleIcon sx={{ mr: 1, my: 0.5 }} />
            <TextField
              id="input:admin.users.modal.username"
              label="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              type="text"
              variant="standard"
              fullWidth
              required
            />
          </FieldRow>

          <FieldRow>
            <LockIcon sx={{ mr: 1, my: 0.5 }} />
            <TextField
              id="input:admin.users.modal.password"
              label="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              type="password"
              variant="standard"
              fullWidth
              required
            />
          </FieldRow>

          <FieldStack id="field:admin.users.modal.roles">
            <FieldLabel variant="text">Roles:</FieldLabel>
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
          </FieldStack>
        </FormStack>
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
