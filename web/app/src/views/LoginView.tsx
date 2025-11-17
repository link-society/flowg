import { useState } from 'react'
import { useNavigate } from 'react-router'

import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import Card from '@mui/material/Card'
import CircularProgress from '@mui/material/CircularProgress'
import Divider from '@mui/material/Divider'
import Grid from '@mui/material/Grid'
import TextField from '@mui/material/TextField'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import LockIcon from '@mui/icons-material/Lock'
import LoginIcon from '@mui/icons-material/Login'

import { UnauthenticatedError } from '@/lib/api/errors'
import * as authApi from '@/lib/api/operations/auth'

import { useApiOperation } from '@/lib/hooks/api'
import { useFeatureFlags } from '@/lib/hooks/featureflags'
import { useNotify } from '@/lib/hooks/notify'

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
    <div className="py-6">
      <header className="mb-6 flex flex-row items-center justify-center gap-2">
        <img src="/web/assets/logo.png" alt="Logo" className="h-16" />
        <h1 className="text-5xl font-bold text-center">FlowG</h1>
      </header>

      <Grid container>
        <Grid size={{ sm: 12, md: 6, lg: 4 }} offset={{ sm: 0, md: 3, lg: 4 }}>
          <Card className="p-3">
            <form
              className="flex flex-col items-stretch gap-3"
              onSubmit={(e) => {
                e.preventDefault()
                handleLogin()
              }}
            >
              <header>
                <h1 className="text-2xl text-center">
                  Sign In with your account
                </h1>
              </header>

              {featureFlags.demoMode && (
                <>
                  <Divider />

                  <p>
                    Demo Mode Enabled, login with <code>demo</code> /{' '}
                    <code>demo</code>.
                  </p>
                </>
              )}

              <Divider />

              <section className="flex flex-col items-stretch gap-3">
                <Box className="flex flex-row items-end">
                  <AccountCircleIcon
                    sx={{ color: 'action.active', mr: 1, my: 0.5 }}
                  />
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
                </Box>

                <Box className="flex flex-row items-end">
                  <LockIcon sx={{ color: 'action.active', mr: 1, my: 0.5 }} />
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
                </Box>
              </section>

              <Divider />

              <Button
                id="btn:login.submit"
                variant="contained"
                color="secondary"
                className="w-full"
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
          </Card>
        </Grid>
      </Grid>
    </div>
  )
}

export default LoginView
