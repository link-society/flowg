import { useState } from 'react'

import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import TextField from '@mui/material/TextField'

import LockIcon from '@mui/icons-material/Lock'
import SendIcon from '@mui/icons-material/Send'

import * as authApi from '@/lib/api/operations/auth'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import { FormRow, IconBox, Label } from './styles'

const PasswordChange = () => {
  const notify = useNotify()

  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')

  const [onSubmit, loading] = useApiOperation(async () => {
    await authApi.changePassword(oldPassword, newPassword)
    notify.success('Password changed')
    setOldPassword('')
    setNewPassword('')
  }, [oldPassword, newPassword, setOldPassword, setNewPassword])

  return (
    <Box>
      <Label variant="text">Change password:</Label>

      <FormRow
        onSubmit={(e) => {
          e.preventDefault()
          onSubmit()
        }}
      >
        <IconBox>
          <LockIcon />
        </IconBox>

        <TextField
          id="input:account.settings.change-password.old"
          label="Old Password"
          value={oldPassword}
          type="password"
          onChange={(e) => setOldPassword(e.target.value)}
          variant="outlined"
          sx={{ flexGrow: 1 }}
          required
        />

        <TextField
          id="input:account.settings.change-password.new"
          label="New Password"
          value={newPassword}
          type="password"
          onChange={(e) => setNewPassword(e.target.value)}
          variant="outlined"
          sx={{ flexGrow: 1 }}
          required
        />

        <Button
          id="btn:account.settings.change-password.submit"
          variant="contained"
          color="secondary"
          sx={{ flexGrow: 0, alignSelf: 'stretch' }}
          type="submit"
          disabled={loading}
        >
          {loading ? (
            <CircularProgress color="inherit" size={24} />
          ) : (
            <SendIcon />
          )}
        </Button>
      </FormRow>
    </Box>
  )
}

export default PasswordChange
