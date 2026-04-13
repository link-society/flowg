import { Typography } from '@mui/material'

import { useState } from 'react'
import { useNavigate } from 'react-router'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import LockIcon from '@mui/icons-material/Lock'
import LoginIcon from '@mui/icons-material/Login'

import { UnauthenticatedError } from '@/lib/api/errors'
import * as authApi from '@/lib/api/operations/auth'

import { useApiOperation } from '@/lib/hooks/api'
import { useFeatureFlags } from '@/lib/hooks/featureflags'
import { useNotify } from '@/lib/hooks/notify'

import {
  LoginViewCard,
  LoginViewCardFields,
  LoginViewContainer,
} from './styles'

const LoginView = () => {
  const featureFlags = useFeatureFlags()
  const navigate = useNavigate()
  const notify = useNotify()

  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')

  const [handleLogin, loading] = useApiOperation(async () => {
    try {
      await authApi.login(username, password)
    } catch (error) {
      if (error instanceof UnauthenticatedError) {
        notify.error('Invalid credentials')
        return
      } else {
        throw error
      }
    }

    navigate('/web/')
  }, [username, password])

  return (
    <LoginViewContainer>
      <header>
        <img src="/web/assets/logo.png" alt="Logo flowG" />
        <Typography variant="h1">FlowG</Typography>
      </header>

      <LoginViewCard>
        <form
          onSubmit={(e) => {
            e.preventDefault()
            handleLogin()
          }}
        >
          <header>
            <Typography variant="h2">Sign In with your account</Typography>
          </header>

          {featureFlags.demoMode && (
            <>
              <Divider />

              <Typography variant="body1">
                Demo Mode Enabled, login with <code>demo</code> /{' '}
                <code>demo</code>.
              </Typography>
            </>
          )}

          <Divider />

          <LoginViewCardFields>
            <div>
              <AccountCircleIcon className="icon" color="action" />
              <TextField
                id="input:login.username"
                label="Username"
                value={username}
                type="text"
                onChange={(e) => setUsername(e.target.value)}
                variant="standard"
                className="grow"
                required
              />
            </div>

            <div>
              <LockIcon className="icon" color="action" />
              <TextField
                id="input:login.password"
                label="Password"
                value={password}
                type="password"
                onChange={(e) => setPassword(e.target.value)}
                variant="standard"
                className="grow"
                disabled={loading}
                required
              />
            </div>
          </LoginViewCardFields>

          <Divider />

          <Button
            id="btn:login.submit"
            variant="contained"
            color="secondary"
            fullWidth
            type="submit"
            startIcon={!loading && <LoginIcon />}
          >
            {loading ? (
              <CircularProgress color="inherit" size={24} />
            ) : (
              <>Sign In</>
            )}
          </Button>
        </form>
      </LoginViewCard>
    </LoginViewContainer>
  )
}

export default LoginView
