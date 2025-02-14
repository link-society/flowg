import React, { useState } from 'react'
import { useApiOperation } from '@/lib/hooks/api'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import LockIcon from '@mui/icons-material/Lock'
import CancelIcon from '@mui/icons-material/Cancel'
import SaveIcon from '@mui/icons-material/Save'

import Dialog from '@mui/material/Dialog'
import DialogTitle from '@mui/material/DialogTitle'
import DialogContent from '@mui/material/DialogContent'
import DialogActions from '@mui/material/DialogActions'
import CircularProgress from '@mui/material/CircularProgress'
import Box from '@mui/material/Box'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import Chip from '@mui/material/Chip'
import Tooltip from '@mui/material/Tooltip'

import { DialogProps } from '@toolpad/core/useDialogs'

import { TransferList } from '@/components/form/transfer-list'

import * as aclApi from '@/lib/api/operations/acls'
import { UserModel } from '@/lib/models'

export const UserFormModal = ({ open, payload, onClose }: DialogProps<string[], UserModel | null>) => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [roles, setRoles] = useState<string[]>([])

  const [onSubmit, loading] = useApiOperation(
    async () => {
      const user = {
        name: username,
        roles,
      }

      await aclApi.saveUser(user, password)
      onClose(user)
    },
    [username, password, roles, onClose],
  )

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
              label="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              type="password"
              variant="standard"
              className="grow"
              required
            />
          </Box>

          <Box className="flex flex-col items-stretch gap-2">
            <span className="font-semibold">Roles:</span>
            <TransferList<string>
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
          variant="contained"
          startIcon={<CancelIcon />}
          onClick={() => onClose(null)}
          disabled={loading}
        >
          Cancel
        </Button>
        <Button
          variant="contained"
          color="secondary"
          startIcon={!loading && <SaveIcon />}
          disabled={loading}
          type="submit"
        >
          {loading
            ? <CircularProgress color="inherit" size={24} />
            : <>Save</>
          }
        </Button>
      </DialogActions>
    </Dialog>
  )
}
