import { useState } from 'react'
import { useTranslation } from 'react-i18next'

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
  const { t } = useTranslation()
  const notify = useNotify()

  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')

  const [onSubmit, loading] = useApiOperation(async () => {
    await authApi.changePassword(oldPassword, newPassword)
    notify.success(t('components.passwordChange.notifications.changed'))
    setOldPassword('')
    setNewPassword('')
  }, [oldPassword, newPassword, setOldPassword, setNewPassword])

  return (
    <Box>
      <Label variant="text">{t('components.passwordChange.title')}</Label>

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
          label={t('components.passwordChange.oldPasswordLabel')}
          value={oldPassword}
          type="password"
          onChange={(e) => setOldPassword(e.target.value)}
          variant="outlined"
          required
        />

        <TextField
          id="input:account.settings.change-password.new"
          label={t('components.passwordChange.newPasswordLabel')}
          value={newPassword}
          type="password"
          onChange={(e) => setNewPassword(e.target.value)}
          variant="outlined"
          required
        />

        <Button
          id="btn:account.settings.change-password.submit"
          variant="contained"
          color="secondary"
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
