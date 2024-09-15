import { useState } from 'react'
import { Form, useNavigate } from 'react-router-dom'
import { useSnackbar } from 'notistack'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import LockIcon from '@mui/icons-material/Lock'
import LoginIcon from '@mui/icons-material/Login'

import Grid from '@mui/material/Grid2'
import Card from '@mui/material/Card'
import Box from '@mui/material/Box'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import { UnauthenticatedError } from '@/lib/api/errors'
import * as authApi from '@/lib/api/operations/auth'

export const LoginView = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)

  const navigate = useNavigate()
  const { enqueueSnackbar } = useSnackbar()

  const handleLogin = async () => {
    setLoading(true)

    try {
      await authApi.login(username, password)
      navigate('/web/')
    }
    catch (error) {
      if (error instanceof UnauthenticatedError) {
        enqueueSnackbar({
          message: 'Invalid credentials',
          variant: 'error'
        })
      }
      else {
        enqueueSnackbar({
          message: 'Unknown error',
          variant: 'error'
        })
      }

      console.error(error)
    }

    setLoading(false)
  }

  return (
    <div className="py-6">
      <Grid container>
        <Grid size={{ sm: 12, md: 6, lg: 4 }} offset={{ sm: 0, md: 3, lg: 4 }}>
          <Card className="p-3">
            <Form
              className="flex flex-col items-stretch gap-3"
              onSubmit={e => {
                e.preventDefault()
                handleLogin()
              }}
            >
              <header>
                <h1 className="text-2xl text-center">Sign In with your account</h1>
              </header>

              <hr />

              <section className="flex flex-col items-stretch gap-3">
                <Box className="flex flex-row items-end">
                  <AccountCircleIcon sx={{ color: 'action.active', mr: 1, my: 0.5 }} />
                  <TextField
                    label="Username"
                    value={username}
                    type="text"
                    onChange={e => setUsername(e.target.value)}
                    variant="standard"
                    className="flex-grow"
                  />
                </Box>

                <Box className="flex flex-row items-end">
                  <LockIcon sx={{ color: 'action.active', mr: 1, my: 0.5 }} />
                  <TextField
                    label="Password"
                    value={password}
                    type="password"
                    onChange={e => setPassword(e.target.value)}
                    variant="standard"
                    className="flex-grow"
                  />
                </Box>
              </section>

              <hr />

              <Button
                variant="contained"
                color="secondary"
                className="w-full"
                type="submit"
                startIcon={!loading && <LoginIcon />}
              >
                {loading
                  ? <CircularProgress color="inherit" />
                  : <>Sign In</>
                }
              </Button>
            </Form>
          </Card>
        </Grid>
      </Grid>
    </div>
  )
}
