import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core'
import { useConfig } from '@/lib/context/config'

import LockIcon from '@mui/icons-material/Lock'
import SendIcon from '@mui/icons-material/Send'

import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import { UnauthenticatedError } from '@/lib/api/errors'

export const PasswordChange = () => {
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [loading, setLoading] = useState(false)

  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()

  const onSubmit = async () => {
    setLoading(true)

    try {
      await new Promise(resolve => setTimeout(resolve, 1000))
      //await userApi.changePassword(oldPassword, newPassword)
      notifications.show('Password changed', {
        severity: 'success',
        autoHideDuration: config.notifications?.autoHideDuration,
      })
    }
    catch (error) {
      if (error instanceof UnauthenticatedError) {
        notifications.show('Session expired', {
          severity: 'error',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
        navigate('/web/login')
      }
      else {
        notifications.show('Unknown error', {
          severity: 'error',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }

      console.error(error)
    }

    setLoading(false)
  }

  return (
    <div>
      <p className="font-semibold mb-2">
        Change password:
      </p>

      <form
        className="flex flex-row items-center gap-2"
        onSubmit={(e) => {
          e.preventDefault()
          onSubmit()
        }}
      >
        <div className="flex-grow-0"><LockIcon /></div>

        <TextField
          label="Old Password"
          value={oldPassword}
          type="password"
          onChange={e => setOldPassword(e.target.value)}
          variant="outlined"
          className="flex-grow"
          required
        />

        <TextField
          label="New Password"
          value={newPassword}
          type="password"
          onChange={e => setNewPassword(e.target.value)}
          variant="outlined"
          className="flex-grow"
          required
        />

        <Button
          variant="contained"
          color="secondary"
          className="flex-grow-0 self-stretch"
          type="submit"
          disabled={loading}
        >
          {loading
            ? <CircularProgress color="inherit" size={24} />
            : <SendIcon />
          }
        </Button>
      </form>
    </div>
  )
}
