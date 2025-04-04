import { useState } from 'react'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import TextField from '@mui/material/TextField'

import LockIcon from '@mui/icons-material/Lock'
import SendIcon from '@mui/icons-material/Send'

import * as authApi from '@/lib/api/operations/auth'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

export const PasswordChange = () => {
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
    <div>
      <p className="font-semibold mb-2">Change password:</p>

      <form
        className="flex flex-row items-center gap-2"
        onSubmit={(e) => {
          e.preventDefault()
          onSubmit()
        }}
      >
        <div className="grow-0">
          <LockIcon />
        </div>

        <TextField
          id="input:account.settings.change-password.old"
          label="Old Password"
          value={oldPassword}
          type="password"
          onChange={(e) => setOldPassword(e.target.value)}
          variant="outlined"
          className="grow"
          required
        />

        <TextField
          id="input:account.settings.change-password.new"
          label="New Password"
          value={newPassword}
          type="password"
          onChange={(e) => setNewPassword(e.target.value)}
          variant="outlined"
          className="grow"
          required
        />

        <Button
          id="btn:account.settings.change-password.submit"
          variant="contained"
          color="secondary"
          className="grow-0 self-stretch"
          type="submit"
          disabled={loading}
        >
          {loading ? (
            <CircularProgress color="inherit" size={24} />
          ) : (
            <SendIcon />
          )}
        </Button>
      </form>
    </div>
  )
}
